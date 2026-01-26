package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"example_shop/service/customer/config"
	"example_shop/service/customer/repository"
)

func main() {
	// 初始化配�?	if err := config.InitConfig(mustResolveConfigPath(
		"service/customer/config/config.yaml",
		"config/config.yaml",
	)); err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	// 初始化数据库
	if err := repository.InitDB(); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	// 迁移数据�?	if err := repository.MigrateTables(); err != nil {
		log.Fatalf("Failed to migrate tables: %v", err)
	}

	if err := installSeedProcedures(); err != nil {
		log.Fatalf("Failed to install seed procedures: %v", err)
	}

	if err := runSeedProcedures(); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	log.Println("Database migration completed successfully!")
}

func installSeedProcedures() error {
	statements := []string{
		`DROP PROCEDURE IF EXISTS sp_seed_shift_config`,
		`CREATE PROCEDURE sp_seed_shift_config()
BEGIN
  DECLARE now_time DATETIME;
  SET now_time = NOW();

  IF NOT EXISTS (SELECT 1 FROM t_shift_config WHERE shift_name='早班' AND is_holiday=0) THEN
    INSERT INTO t_shift_config (shift_name,start_time,end_time,min_staff,is_holiday,create_time,update_time,create_by)
    VALUES ('早班',TIMESTAMP('2000-01-01','08:00:00'),TIMESTAMP('2000-01-01','16:00:00'),5,0,now_time,now_time,'ADMIN_SEED');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM t_shift_config WHERE shift_name='中班' AND is_holiday=0) THEN
    INSERT INTO t_shift_config (shift_name,start_time,end_time,min_staff,is_holiday,create_time,update_time,create_by)
    VALUES ('中班',TIMESTAMP('2000-01-01','12:00:00'),TIMESTAMP('2000-01-01','20:00:00'),5,0,now_time,now_time,'ADMIN_SEED');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM t_shift_config WHERE shift_name='晚班' AND is_holiday=0) THEN
    INSERT INTO t_shift_config (shift_name,start_time,end_time,min_staff,is_holiday,create_time,update_time,create_by)
    VALUES ('晚班',TIMESTAMP('2000-01-01','16:00:00'),TIMESTAMP('2000-01-01','00:00:00'),3,0,now_time,now_time,'ADMIN_SEED');
  END IF;

  IF NOT EXISTS (SELECT 1 FROM t_shift_config WHERE shift_name='夜班' AND is_holiday=0) THEN
    INSERT INTO t_shift_config (shift_name,start_time,end_time,min_staff,is_holiday,create_time,update_time,create_by)
    VALUES ('夜班',TIMESTAMP('2000-01-01','00:00:00'),TIMESTAMP('2000-01-01','08:00:00'),2,0,now_time,now_time,'ADMIN_SEED');
  END IF;
END`,

		`DROP PROCEDURE IF EXISTS sp_seed_customer_service`,
		`CREATE PROCEDURE sp_seed_customer_service()
BEGIN
  DECLARE now_time DATETIME;
  SET now_time = NOW();

  INSERT INTO t_customer_service (cs_id,cs_name,dept_id,team_id,skill_tags,status,current_status,create_time,update_time)
  SELECT 'KF001','张三','DEPT001','TEAM001','票务查询,退款处�?,1,0,now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_customer_service WHERE cs_id='KF001');

  INSERT INTO t_customer_service (cs_id,cs_name,dept_id,team_id,skill_tags,status,current_status,create_time,update_time)
  SELECT 'KF002','李四','DEPT001','TEAM001','改签咨询,订单查询',1,0,now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_customer_service WHERE cs_id='KF002');

  INSERT INTO t_customer_service (cs_id,cs_name,dept_id,team_id,skill_tags,status,current_status,create_time,update_time)
  SELECT 'KF003','王五','DEPT001','TEAM002','投诉处理,建议收集',1,1,now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_customer_service WHERE cs_id='KF003');

  INSERT INTO t_customer_service (cs_id,cs_name,dept_id,team_id,skill_tags,status,current_status,create_time,update_time)
  SELECT 'KF004','赵六','DEPT002','TEAM003','演出票务,座位咨询',1,0,now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_customer_service WHERE cs_id='KF004');

  INSERT INTO t_customer_service (cs_id,cs_name,dept_id,team_id,skill_tags,status,current_status,create_time,update_time)
  SELECT 'KF005','钱七','DEPT002','TEAM003','发票开�?支付问题',1,0,now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_customer_service WHERE cs_id='KF005');

  INSERT INTO t_customer_service (cs_id,cs_name,dept_id,team_id,skill_tags,status,current_status,create_time,update_time)
  SELECT 'KF006','孙八','DEPT003','TEAM004','技术支�?账号问题',1,2,now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_customer_service WHERE cs_id='KF006');
END`,

		`DROP PROCEDURE IF EXISTS sp_seed_quick_reply`,
		`CREATE PROCEDURE sp_seed_quick_reply()
BEGIN
  DECLARE now_time DATETIME;
  SET now_time = NOW();

  DELETE FROM t_quick_reply WHERE reply_content LIKE '【种子�?';

  INSERT INTO t_quick_reply (reply_type,reply_content,create_by,is_public,create_time,update_time)
  VALUES
    (3,'【种子】您好，欢迎咨询天极票务，请问需要我帮您查询哪个订单�?,'KF001',1,now_time,now_time),
    (0,'【种子】请提供订单�?手机号后四位，我马上为您核实�?,'KF001',1,now_time,now_time),
    (0,'【种子】退款一�?-7个工作日原路退回，请您耐心等待�?,'KF002',1,now_time,now_time),
    (1,'【种子】非常抱歉给您带来不便，我先为您登记并立即跟进�?,'KF003',1,now_time,now_time),
    (2,'【种子】感谢您的建议，我们会记录并持续优化体验�?,'KF003',1,now_time,now_time),
    (3,'【种子】已为您处理完成，如还有问题随时联系我�?,'KF002',1,now_time,now_time);
END`,

		`DROP PROCEDURE IF EXISTS sp_seed_schedule`,
		`CREATE PROCEDURE sp_seed_schedule(IN p_start_date DATE, IN p_days INT)
BEGIN
  DECLARE i INT DEFAULT 0;
  DECLARE d DATE;
  DECLARE now_time DATETIME;
  SET now_time = NOW();

  IF p_days IS NULL OR p_days <= 0 THEN
    SET p_days = 7;
  END IF;

  DELETE FROM t_schedule
  WHERE cs_id IN ('KF001','KF002','KF003','KF004','KF005','KF006')
    AND schedule_date BETWEEN p_start_date AND DATE_ADD(p_start_date, INTERVAL p_days - 1 DAY);

  WHILE i < p_days DO
    SET d = DATE_ADD(p_start_date, INTERVAL i DAY);

    INSERT INTO t_schedule (cs_id,shift_id,schedule_date,status,replace_cs_id,create_time,update_time)
    SELECT 'KF001', shift_id, d, 0, NULL, now_time, now_time FROM t_shift_config WHERE shift_name='早班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;

    INSERT INTO t_schedule (cs_id,shift_id,schedule_date,status,replace_cs_id,create_time,update_time)
    SELECT 'KF002', shift_id, d, 0, NULL, now_time, now_time FROM t_shift_config WHERE shift_name='中班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;

    INSERT INTO t_schedule (cs_id,shift_id,schedule_date,status,replace_cs_id,create_time,update_time)
    SELECT 'KF004', shift_id, d, 0, NULL, now_time, now_time FROM t_shift_config WHERE shift_name='晚班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;

    INSERT INTO t_schedule (cs_id,shift_id,schedule_date,status,replace_cs_id,create_time,update_time)
    SELECT 'KF005', shift_id, d, 0, NULL, now_time, now_time FROM t_shift_config WHERE shift_name='夜班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;

    IF MOD(i, 5) = 2 THEN
      INSERT INTO t_schedule (cs_id,shift_id,schedule_date,status,replace_cs_id,create_time,update_time)
      SELECT 'KF003', shift_id, d, 1, NULL, now_time, now_time FROM t_shift_config WHERE shift_name='中班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;
    ELSE
      INSERT INTO t_schedule (cs_id,shift_id,schedule_date,status,replace_cs_id,create_time,update_time)
      SELECT 'KF003', shift_id, d, 0, NULL, now_time, now_time FROM t_shift_config WHERE shift_name='中班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;
    END IF;

    IF MOD(i, 7) = 3 THEN
      INSERT INTO t_schedule (cs_id,shift_id,schedule_date,status,replace_cs_id,create_time,update_time)
      SELECT 'KF006', shift_id, d, 2, 'KF001', now_time, now_time FROM t_shift_config WHERE shift_name='早班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;
    ELSE
      INSERT INTO t_schedule (cs_id,shift_id,schedule_date,status,replace_cs_id,create_time,update_time)
      SELECT 'KF006', shift_id, d, 0, NULL, now_time, now_time FROM t_shift_config WHERE shift_name='早班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;
    END IF;

    SET i = i + 1;
  END WHILE;
END`,

		`DROP PROCEDURE IF EXISTS sp_seed_conversation`,
		`CREATE PROCEDURE sp_seed_conversation(IN p_conv_count INT)
BEGIN
  DECLARE i INT DEFAULT 1;
  DECLARE now_time DATETIME;
  DECLARE conv_id VARCHAR(64);
  DECLARE start_t DATETIME;
  DECLARE end_t DATETIME;
  DECLARE cs_id VARCHAR(32);
  DECLARE src VARCHAR(32);
  DECLARE msg_type_val TINYINT;
  SET now_time = NOW();

  IF p_conv_count IS NULL OR p_conv_count <= 0 THEN
    SET p_conv_count = 6;
  END IF;

  DELETE m FROM t_conv_message m INNER JOIN t_conversation c ON m.conv_id=c.conv_id WHERE c.conv_id LIKE 'SEED-%';
  DELETE FROM t_conversation WHERE conv_id LIKE 'SEED-%';

  WHILE i <= p_conv_count DO
    SET conv_id = CONCAT('SEED-', DATE_FORMAT(now_time, '%Y%m%d%H%i%s'), '-', LPAD(i, 3, '0'));
    SET start_t = DATE_SUB(now_time, INTERVAL (p_conv_count - i) HOUR);
    SET end_t = DATE_ADD(start_t, INTERVAL 12 MINUTE);

    IF MOD(i, 3) = 1 THEN
      SET cs_id = 'KF001';
      SET src = 'APP';
      SET msg_type_val = 0;
    ELSEIF MOD(i, 3) = 2 THEN
      SET cs_id = 'KF003';
      SET src = '网页';
      SET msg_type_val = 1;
    ELSE
      SET cs_id = 'KF005';
      SET src = 'H5';
      SET msg_type_val = 0;
    END IF;

    INSERT INTO t_conversation
      (conv_id,user_id,user_nickname,cs_id,transfer_cs_id,source,start_time,end_time,msg_type,is_manual_adjust,status,create_time,update_time)
    VALUES
      (conv_id, CONCAT('U', LPAD(i, 4, '0')), CONCAT('用户', i), cs_id, NULL, src, start_t, end_t, msg_type_val, 0, 1, start_t, end_t);

    INSERT INTO t_conv_message
      (conv_id,sender_type,sender_id,msg_content,file_url,file_type,voice_url,is_quick_reply,quick_reply_id,send_time)
    VALUES
      (conv_id,0,CONCAT('U', LPAD(i, 4, '0')),'我想查询订单状态，订单号是TJ20260113-00',NULL,NULL,NULL,0,NULL,DATE_ADD(start_t, INTERVAL 1 MINUTE));

    INSERT INTO t_conv_message
      (conv_id,sender_type,sender_id,msg_content,file_url,file_type,voice_url,is_quick_reply,quick_reply_id,send_time)
    SELECT
      conv_id,1,cs_id,qr.reply_content,NULL,NULL,NULL,1,qr.reply_id,DATE_ADD(start_t, INTERVAL 3 MINUTE)
    FROM t_quick_reply qr
    WHERE qr.reply_content LIKE '【种子】请提供订单�?'
    ORDER BY qr.reply_id DESC LIMIT 1;

    INSERT INTO t_conv_message
      (conv_id,sender_type,sender_id,msg_content,file_url,file_type,voice_url,is_quick_reply,quick_reply_id,send_time)
    VALUES
      (conv_id,0,CONCAT('U', LPAD(i, 4, '0')),'好的，手机号后四位是8899�?,NULL,NULL,NULL,0,NULL,DATE_ADD(start_t, INTERVAL 5 MINUTE));

    SET i = i + 1;
  END WHILE;
END`,

		`DROP PROCEDURE IF EXISTS sp_seed_leave_transfer`,
		`CREATE PROCEDURE sp_seed_leave_transfer(IN p_start_date DATE)
BEGIN
  DECLARE now_time DATETIME;
  DECLARE day2 DATE;
  DECLARE day4 DATE;
  DECLARE shift_early BIGINT;
  DECLARE shift_mid BIGINT;
  SET now_time = NOW();
  SET day2 = DATE_ADD(p_start_date, INTERVAL 2 DAY);
  SET day4 = DATE_ADD(p_start_date, INTERVAL 4 DAY);

  SELECT shift_id INTO shift_early FROM t_shift_config WHERE shift_name='早班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;
  SELECT shift_id INTO shift_mid FROM t_shift_config WHERE shift_name='中班' AND is_holiday=0 ORDER BY shift_id LIMIT 1;

  DELETE FROM t_leave_transfer WHERE reason LIKE '【种子�?';

  INSERT INTO t_leave_transfer
    (cs_id,apply_type,target_date,shift_id,target_cs_id,approval_status,approver_id,approval_time,reason,create_time,update_time)
  VALUES
    ('KF006',0,day2,shift_early,NULL,1,'ADMIN001',now_time,'【种子】身体不适，请假一�?,now_time,now_time),
    ('KF003',0,day4,shift_mid,NULL,0,NULL,NULL,'【种子】家中有事，申请请假',now_time,now_time),
    ('KF004',1,day2,shift_mid,'KF005',1,'ADMIN001',now_time,'【种子】申请调班至KF005',now_time,now_time);
END`,

		`DROP PROCEDURE IF EXISTS sp_seed_conv_tag`,
		`CREATE PROCEDURE sp_seed_conv_tag()
BEGIN
  DECLARE now_time DATETIME;
  SET now_time = NOW();

  INSERT INTO t_conv_tag (tag_name,tag_color,sort_no,create_by,create_time,update_time)
  SELECT '投诉','#f5222d',1,'ADMIN_SEED',now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_conv_tag WHERE tag_name='投诉');

  INSERT INTO t_conv_tag (tag_name,tag_color,sort_no,create_by,create_time,update_time)
  SELECT '退�?,'#fa8c16',2,'ADMIN_SEED',now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_conv_tag WHERE tag_name='退�?);

  INSERT INTO t_conv_tag (tag_name,tag_color,sort_no,create_by,create_time,update_time)
  SELECT '改签','#1890ff',3,'ADMIN_SEED',now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_conv_tag WHERE tag_name='改签');

  INSERT INTO t_conv_tag (tag_name,tag_color,sort_no,create_by,create_time,update_time)
  SELECT '支付','#52c41a',4,'ADMIN_SEED',now_time,now_time FROM DUAL
  WHERE NOT EXISTS (SELECT 1 FROM t_conv_tag WHERE tag_name='支付');
END`,

		`DROP PROCEDURE IF EXISTS sp_seed_all`,
		`CREATE PROCEDURE sp_seed_all(IN p_start_date DATE, IN p_days INT, IN p_conv_count INT)
BEGIN
  CALL sp_seed_shift_config();
  CALL sp_seed_customer_service();
  CALL sp_seed_quick_reply();
  CALL sp_seed_schedule(p_start_date, p_days);
  CALL sp_seed_conversation(p_conv_count);
  CALL sp_seed_leave_transfer(p_start_date);
  CALL sp_seed_conv_tag();
END`,
	}

	for _, sql := range statements {
		if err := repository.DB.Exec(sql).Error; err != nil {
			return fmt.Errorf("exec failed: %w; sql=%s", err, shortSQL(sql))
		}
	}
	return nil
}

func runSeedProcedures() error {
	if err := repository.DB.Exec("CALL sp_seed_all(CURDATE(), ?, ?)", 7, 8).Error; err != nil {
		return err
	}
	return nil
}

func shortSQL(sql string) string {
	const maxLen = 160
	if len(sql) <= maxLen {
		return sql
	}
	return sql[:maxLen] + "..."
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
