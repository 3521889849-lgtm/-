package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"example_shop/service/customer/config"
	"example_shop/service/customer/repository"
	"example_shop/service/customer/model"

	"gorm.io/gorm"
)

func main() {
	if err := config.InitConfig(mustResolveConfigPath(
		"service/customer/config/config.yaml",
		"config/config.yaml",
	)); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	if err := repository.MigrateTables(); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	if err := seedHistoryConversations(repository.DB, 10); err != nil {
		log.Fatalf("Failed to seed history conversations: %v", err)
	}

	log.Println("Seed history conversations completed successfully!")
}

func seedHistoryConversations(db *gorm.DB, n int) error {
	if n <= 0 {
		n = 10
	}

	now := time.Now()
	prefix := "HSEED-" + now.Format("20060102150405")

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("conv_id LIKE ?", "HSEED-%").Delete(&model.ConvMessage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("conv_id LIKE ?", "HSEED-%").Delete(&model.Conversation{}).Error; err != nil {
			return err
		}

		csIDs := []string{"KF001", "KF002", "KF003"}
		sources := []string{"网页", "APP", "H5"}

		for i := 1; i <= n; i++ {
			convID := fmt.Sprintf("%s-%03d", prefix, i)
			userID := fmt.Sprintf("U9%03d", i)
			userNick := fmt.Sprintf("模拟用户%d", i)
			csID := csIDs[(i-1)%len(csIDs)]
			source := sources[(i-1)%len(sources)]

			startTime := now.Add(-time.Duration(n-i+1) * time.Hour)
			endTime := startTime.Add(8 * time.Minute)

			status := int8(1) // 所有历史会话都标记为已结束

			conv := model.Conversation{
				ConvID:         convID,
				UserID:         userID,
				UserNickname:   userNick,
				CsID:           csID,
				Source:         source,
				StartTime:      startTime,
				EndTime:        endTime,
				MsgType:        0,
				IsManualAdjust: 0,
				CategoryID:     0,
				Tags:           "模拟数据",
				IsCore:         0,
				Status:         status,
				CreateTime:     startTime,
				UpdateTime:     endTime,
			}
			if err := tx.Create(&conv).Error; err != nil {
				return err
			}

			msgs := []model.ConvMessage{
				{
					ConvID:       convID,
					SenderType:   0,
					SenderID:     userID,
					MsgContent:   fmt.Sprintf("你好，我想查询订单状态（测试会话 %d）�?, i),
					IsQuickReply: 0,
					SendTime:     startTime.Add(1 * time.Minute),
				},
				{
					ConvID:       convID,
					SenderType:   2,
					SenderID:     "SYSTEM",
					MsgContent:   "【系统】会话已结束，归档入历史记录�?,
					IsQuickReply: 0,
					SendTime:     startTime.Add(2 * time.Minute),
				},
				{
					ConvID:       convID,
					SenderType:   1,
					SenderID:     csID,
					MsgContent:   "您好，我已为您核实订单状态，请稍等�?,
					IsQuickReply: 0,
					SendTime:     startTime.Add(3 * time.Minute),
				},
				{
					ConvID:       convID,
					SenderType:   0,
					SenderID:     userID,
					MsgContent:   "好的，谢谢�?,
					IsQuickReply: 0,
					SendTime:     startTime.Add(4 * time.Minute),
				},
			}
			if err := tx.Create(&msgs).Error; err != nil {
				return err
			}

			if err := tx.Model(&model.Conversation{}).Where("conv_id = ?", convID).Update("update_time", endTime).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func mustResolveConfigPath(candidates ...string) string {
	for _, p := range candidates {
		if fileExists(p) {
			return p
		}
	}

	root, ok := findProjectRoot()
	if ok {
		for _, p := range candidates {
			if fileExists(filepath.Join(root, p)) {
				return filepath.Join(root, p)
			}
		}
	}

	return candidates[0]
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func findProjectRoot() (string, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	dir := wd
	for {
		if fileExists(filepath.Join(dir, "go.mod")) {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
