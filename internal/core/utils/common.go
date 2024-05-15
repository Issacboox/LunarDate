package utils

import (
	"os"
	"strconv"
)

func GetEnv(key, fallback string) string {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	// จัดการข้อผิดพลาดที่เกิดจากไม่สามารถอ่านไฟล์ .env ได้
	// 	// ยกตัวอย่างเช่นการแสดงข้อความแจ้งเตือนหรือใช้ค่าเริ่มต้นที่คุณเลือก
	// 	// fallback หรือค่าว่างเปล่าเป็นตัวเลือกของคุณ
	// 	return fallback
	// }
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func AtoI(s string, v int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return v
	}
	return i
}