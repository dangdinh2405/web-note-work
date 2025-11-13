package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dangdinh2405/web-note-work/internal/models"
)

func GetAllTasks(taskColl *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := c.DefaultQuery("filter", "today") 

		now := time.Now()
		var startDate time.Time

		switch filter {
		case "today":
			startDate = time.Date(
				now.Year(), now.Month(), now.Day(),
				0, 0, 0, 0, now.Location(),
			)

		case "week":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			monday := now.AddDate(0, 0, -(weekday-1))
			startDate = time.Date(
				monday.Year(), monday.Month(), monday.Day(),
				0, 0, 0, 0, now.Location(),
			)

		case "month":
			startDate = time.Date(
				now.Year(), now.Month(), 1,
				0, 0, 0, 0, now.Location(),
			)

		case "all":
			fallthrough
		default:
			startDate = time.Time{}
		}

		matchCond := bson.D{}
		if !startDate.IsZero() {
			matchCond = bson.D{
				{"createdAt", bson.D{
					{"$gte", startDate},
				}},
			}
		}

		pipeline := mongo.Pipeline{
			{
				{"$match", matchCond},
			},
			{
				{"$facet", bson.D{
					{"tasks", mongo.Pipeline{
						{{"$sort", bson.D{{"createdAt", -1}}}},
					}},
					{"activeCount", mongo.Pipeline{
						{{"$match", bson.D{{"status", "active"}}}},
						{{"$count", "count"}},
					}},
					{"completeCount", mongo.Pipeline{
						{{"$match", bson.D{{"status", "complete"}}}},
						{{"$count", "count"}},
					}},
				}},
			},
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		cursor, err := taskColl.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi hệ thống"})
			return
		}
		defer cursor.Close(ctx)

		type facetResult struct {
			Tasks         []models.Task `bson:"tasks" json:"tasks"`
			ActiveCount   []struct {
				Count int64 `bson:"count"`
			} `bson:"activeCount"`
			CompleteCount []struct {
				Count int64 `bson:"count"`
			} `bson:"completeCount"`
		}

		var results []facetResult
		if err := cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi hệ thống"})
			return
		}

		if len(results) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"tasks":         []models.Task{},
				"activeCount":   0,
				"completeCount": 0,
			})
			return
		}

		res := results[0]

		var activeCount int64
		if len(res.ActiveCount) > 0 {
			activeCount = res.ActiveCount[0].Count
		}

		var completeCount int64
		if len(res.CompleteCount) > 0 {
			completeCount = res.CompleteCount[0].Count
		}

		c.JSON(http.StatusOK, gin.H{
			"tasks":         res.Tasks,
			"activeCount":   activeCount,
			"completeCount": completeCount,
		})
	}
}

func CreateTask(taskColl *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Title string `json:"title"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Dữ liệu không hợp lệ"})
			return
		}

		now := time.Now()
		task := models.Task{
			Title:     body.Title,
			Status:    "active",
			CreatedAt: now,
			UpdatedAt: now,
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		result, err := taskColl.InsertOne(ctx, task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi hệ thống"})
			return
		}

		taskWithID := gin.H{
			"_id":        result.InsertedID,
			"title":      task.Title,
			"status":     task.Status,
			"createdAt":  task.CreatedAt,
			"updatedAt":  task.UpdatedAt,
			"completedAt": task.CompletedAt,
		}

		c.JSON(http.StatusCreated, taskWithID)
	}
}

func UpdateTask(taskColl *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ID không hợp lệ"})
			return
		}

		var body struct {
			Title       *string `json:"title"`
			Status      *string `json:"status"`
			CompletedAt *string `json:"completedAt"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Dữ liệu không hợp lệ"})
			return
		}

		updateFields := bson.M{}

		if body.Title != nil {
			updateFields["title"] = *body.Title
		}

		if body.Status != nil {
			updateFields["status"] = *body.Status
		}

		if body.CompletedAt != nil {
			if *body.CompletedAt == "" {
				updateFields["completedAt"] = nil
			} else {
				t, err := time.Parse(time.RFC3339, *body.CompletedAt)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": "completedAt không đúng định dạng ISO"})
					return
				}
				updateFields["completedAt"] = t
			}
		}


		if len(updateFields) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Không có field nào để cập nhật"})
			return
		}
		updateFields["updatedAt"] = time.Now()

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		opts := options.FindOneAndUpdate().
			SetReturnDocument(options.After)

		var updatedTask models.Task

		err = taskColl.FindOneAndUpdate(
			ctx,
			bson.M{"_id": objID},
			bson.M{"$set": updateFields},
			opts,
		).Decode(&updatedTask)

		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"message": "Nhiệm vụ không tồn tại"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi hệ thống"})
			return
		}

		c.JSON(http.StatusOK, updatedTask)
	}
}

func DeleteTask(taskColl *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "ID không hợp lệ"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		var deletedTask models.Task

		err = taskColl.FindOneAndDelete(ctx, bson.M{"_id": objID}).Decode(&deletedTask)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"message": "Nhiệm vụ không tồn tại"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi hệ thống"})
			return
		}

		c.JSON(http.StatusOK, deletedTask)
	}
}