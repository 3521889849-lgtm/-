package service

import (
	"errors"
	"example_shop/common/db"
	"example_shop/common/model/hotel_admin"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nfnt/resize"
	"gorm.io/gorm"
)

// RoomImageService 房源图片管理服务
// 负责处理房源图片的上传、删除等核心业务逻辑，
// 包括图片数量限制、图片格式验证、图片尺寸压缩、图片存储路径管理等。
type RoomImageService struct{}

const (
	MaxImageCount  = 16                    // 最大图片数量
	ImageWidth     = 400                   // 图片宽度
	ImageHeight    = 300                   // 图片高度
	MaxFileSize    = 5 * 1024 * 1024       // 最大文件大小 5MB
	UploadBasePath = "uploads/room_images" // 上传基础路径
)

// UploadRoomImages 批量上传房源图片
func (s *RoomImageService) UploadRoomImages(roomID uint64, files []*multipart.FileHeader) ([]hotel_admin.RoomImage, error) {
	// 检查房源是否存在
	var roomInfo hotel_admin.RoomInfo                                 // 声明房源实体变量，用于存储查询到的房源信息
	if err := db.MysqlDB.First(&roomInfo, roomID).Error; err != nil { // 通过房源ID查询房源信息，如果查询失败则说明房源不存在
		return nil, errors.New("房源不存在") // 返回nil和错误信息，表示房源不存在
	}

	// 检查当前图片数量
	var currentCount int64                                                                                              // 声明计数变量，用于存储当前房源的图片数量
	db.MysqlDB.Model(&hotel_admin.RoomImage{}).Where("room_id = ? AND deleted_at IS NULL", roomID).Count(&currentCount) // 统计该房源当前的图片数量（排除已删除图片）

	// 检查上传数量限制
	if int(currentCount)+len(files) > MaxImageCount { // 如果当前图片数量加上待上传文件数量大于最大图片数量限制
		return nil, fmt.Errorf("图片数量超过限制，最多只能上传 %d 张，当前已有 %d 张", MaxImageCount, currentCount) // 返回nil和错误信息，表示图片数量超过限制（包含最大数量和当前数量）
	}

	// 创建上传目录
	if err := os.MkdirAll(UploadBasePath, 0755); err != nil { // 创建上传目录（如果不存在），设置目录权限为0755（所有者可读写执行，组和其他用户可读执行），如果创建失败则返回错误
		return nil, fmt.Errorf("创建上传目录失败: %v", err) // 返回nil和错误信息，表示创建上传目录失败
	}

	var uploadedImages []hotel_admin.RoomImage // 声明已上传图片列表变量，用于存储成功上传的图片信息
	nextSortOrder := uint8(currentCount)       // 计算下一个排序序号（从当前图片数量开始，确保新上传的图片排在最后）

	for i, fileHeader := range files { // 遍历待上传的文件列表
		// 验证文件大小
		if fileHeader.Size > MaxFileSize { // 如果文件大小超过最大文件大小限制（5MB）
			return nil, fmt.Errorf("文件 %s 大小超过限制（最大 5MB）", fileHeader.Filename) // 返回nil和错误信息，表示文件大小超过限制（包含文件名）
		}

		// 验证文件格式
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename)) // 获取文件扩展名并转换为小写（用于格式验证）
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {     // 如果文件扩展名不是jpg、jpeg或png格式
			return nil, fmt.Errorf("文件 %s 格式不支持，仅支持 jpg/png 格式", fileHeader.Filename) // 返回nil和错误信息，表示文件格式不支持（包含文件名）
		}

		// 打开文件
		file, err := fileHeader.Open() // 打开文件（获取文件读取器），如果打开失败则返回错误
		if err != nil {                // 如果打开文件失败，则返回错误
			return nil, fmt.Errorf("打开文件失败: %v", err) // 返回nil和错误信息，表示打开文件失败
		}
		defer file.Close() // 延迟关闭文件（确保资源释放）

		// 解码图片
		var img image.Image // 声明图片对象变量，用于存储解码后的图片数据
		var format string   // 声明格式变量，用于存储图片格式（png或jpg）
		if ext == ".png" {  // 如果文件扩展名是png格式
			img, err = png.Decode(file) // 使用PNG解码器解码图片，如果解码失败则返回错误
			format = "png"              // 设置格式为png
		} else {
			img, err = jpeg.Decode(file) // 使用JPEG解码器解码图片，如果解码失败则返回错误
			format = "jpg"               // 设置格式为jpg
		}
		if err != nil { // 如果图片解码失败，则返回错误
			return nil, fmt.Errorf("图片解码失败: %v", err) // 返回nil和错误信息，表示图片解码失败
		}

		// 重置文件指针
		file.Seek(0, 0) // 将文件指针重置到文件开头（用于后续可能的重新读取）

		// 调整图片尺寸
		resizedImg := resize.Resize(ImageWidth, ImageHeight, img, resize.Lanczos3) // 调整图片尺寸到指定宽度和高度（使用Lanczos3算法，保证图片质量）

		// 生成文件名
		fileName := fmt.Sprintf("%d_%d_%d.%s", roomID, time.Now().Unix(), i, format) // 生成唯一文件名（格式：房源ID_时间戳_索引.格式）
		filePath := filepath.Join(UploadBasePath, fileName)                          // 构建完整文件路径（拼接上传基础路径和文件名）

		// 保存调整后的图片
		outFile, err := os.Create(filePath) // 创建输出文件（用于保存调整后的图片），如果创建失败则返回错误
		if err != nil {                     // 如果创建文件失败，则返回错误
			return nil, fmt.Errorf("创建文件失败: %v", err) // 返回nil和错误信息，表示创建文件失败
		}
		defer outFile.Close() // 延迟关闭输出文件（确保资源释放）

		if format == "png" { // 如果图片格式是png
			err = png.Encode(outFile, resizedImg) // 使用PNG编码器将调整后的图片编码并写入输出文件
		} else {
			err = jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 90}) // 使用JPEG编码器将调整后的图片编码并写入输出文件（设置JPEG质量为90）
		}
		if err != nil { // 如果图片编码或写入失败，则返回错误
			os.Remove(filePath)                       // 删除已创建的文件（清理失败的文件）
			return nil, fmt.Errorf("保存图片失败: %v", err) // 返回nil和错误信息，表示保存图片失败
		}

		// 生成访问URL（实际项目中应该使用CDN或对象存储）
		imageURL := fmt.Sprintf("/%s/%s", UploadBasePath, fileName) // 生成图片访问URL（格式：/上传基础路径/文件名）

		// 保存到数据库
		roomImage := hotel_admin.RoomImage{ // 创建房源图片实体对象
			RoomID:      roomID,                                        // 设置房源ID（从参数中获取）
			ImageURL:    imageURL,                                      // 设置图片URL（生成的访问URL）
			ImageSize:   fmt.Sprintf("%dx%d", ImageWidth, ImageHeight), // 设置图片尺寸（格式：宽度x高度）
			ImageFormat: format,                                        // 设置图片格式（png或jpg）
			SortOrder:   nextSortOrder + uint8(i),                      // 设置排序序号（下一个排序序号加上当前索引，确保图片按上传顺序排列）
		}

		if err := db.MysqlDB.Create(&roomImage).Error; err != nil { // 将图片信息保存到数据库，如果保存失败则返回错误
			os.Remove(filePath)                         // 删除已保存的文件（清理失败的文件）
			return nil, fmt.Errorf("保存图片记录失败: %v", err) // 返回nil和错误信息，表示保存图片记录失败
		}

		uploadedImages = append(uploadedImages, roomImage) // 将成功上传的图片信息添加到已上传图片列表中
	}

	return uploadedImages, nil // 返回已上传图片列表和无错误
}

// DeleteRoomImage 删除房源图片
func (s *RoomImageService) DeleteRoomImage(imageID uint64) error {
	var roomImage hotel_admin.RoomImage                                 // 声明房源图片实体变量，用于存储查询到的图片信息
	if err := db.MysqlDB.First(&roomImage, imageID).Error; err != nil { // 通过图片ID查询图片信息，如果查询失败则说明图片不存在
		return errors.New("图片不存在") // 返回错误信息，表示图片不存在
	}

	// 删除文件
	if strings.HasPrefix(roomImage.ImageURL, "/") { // 如果图片URL以"/"开头（本地文件路径格式）
		filePath := strings.TrimPrefix(roomImage.ImageURL, "/") // 移除URL开头的"/"，获取实际文件路径
		if _, err := os.Stat(filePath); err == nil {            // 检查文件是否存在（如果文件存在，则err为nil）
			os.Remove(filePath) // 删除物理文件（从文件系统中删除）
		}
	}

	// 软删除数据库记录
	return db.MysqlDB.Delete(&roomImage).Error // 执行软删除操作（设置deleted_at字段），根据图片ID删除图片记录，返回删除操作的结果（成功为nil，失败为error）
}

// GetRoomImages 获取房源图片列表
func (s *RoomImageService) GetRoomImages(roomID uint64) ([]hotel_admin.RoomImage, error) {
	var images []hotel_admin.RoomImage                                        // 声明图片列表变量，用于存储查询到的图片信息列表
	if err := db.MysqlDB.Where("room_id = ? AND deleted_at IS NULL", roomID). // 添加筛选条件：房源ID匹配且未被删除
											Order("sort_order ASC").          // 添加排序条件，按排序序号正序排列（排序序号小的图片排在前面）
											Find(&images).Error; err != nil { // 执行查询并获取符合条件的图片列表，如果查询失败则返回错误
		return nil, err // 返回nil和错误信息
	}
	return images, nil // 返回图片列表和无错误
}

// UpdateImageSortOrder 更新图片排序
func (s *RoomImageService) UpdateImageSortOrder(imageID uint64, sortOrder uint8) error {
	var roomImage hotel_admin.RoomImage                                 // 声明房源图片实体变量，用于存储查询到的图片信息
	if err := db.MysqlDB.First(&roomImage, imageID).Error; err != nil { // 通过图片ID查询图片信息，如果查询失败则说明图片不存在
		return errors.New("图片不存在") // 返回错误信息，表示图片不存在
	}

	roomImage.SortOrder = sortOrder          // 更新排序序号（从参数中获取）
	return db.MysqlDB.Save(&roomImage).Error // 保存图片信息到数据库，返回保存操作的结果（成功为nil，失败为error）
}

// BatchUpdateImageSortOrder 批量更新图片排序
func (s *RoomImageService) BatchUpdateImageSortOrder(roomID uint64, imageSortOrders []ImageSortOrder) error {
	if len(imageSortOrders) == 0 { // 如果排序列表为空（长度为0）
		return errors.New("排序列表不能为空") // 返回错误信息，表示排序列表不能为空
	}

	// 使用事务批量更新
	return db.MysqlDB.Transaction(func(tx *gorm.DB) error { // 开启数据库事务，传入事务处理函数，返回事务执行结果
		for _, item := range imageSortOrders { // 遍历排序列表，为每个图片更新排序序号
			if err := tx.Model(&hotel_admin.RoomImage{}). // 创建房源图片模型的查询构建器（使用事务连接）
									Where("id = ? AND room_id = ?", item.ImageID, roomID).   // 添加筛选条件：图片ID匹配且房源ID匹配（确保只更新指定房源的图片）
									Update("sort_order", item.SortOrder).Error; err != nil { // 更新排序序号字段，如果更新失败则返回错误
				return err // 返回数据库操作错误
			}
		}
		return nil // 返回nil表示事务执行成功（所有图片的排序序号都更新成功）
	})
}

// ImageSortOrder 图片排序
type ImageSortOrder struct {
	ImageID   uint64 `json:"image_id"`
	SortOrder uint8  `json:"sort_order"`
}
