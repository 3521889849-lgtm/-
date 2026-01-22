DROP PROCEDURE IF EXISTS InitMockData;

DELIMITER //

CREATE PROCEDURE InitMockData()
BEGIN
    -- 1. 初始化班次配置 (t_shift_config)
    DELETE FROM t_shift_config;
    INSERT INTO t_shift_config (shift_name, start_time, end_time, min_staff, is_holiday, create_time, update_time, create_by) VALUES 
    ('早班', '09:00:00', '18:00:00', 3, 0, NOW(), NOW(), 'system'),
    ('中班', '13:00:00', '22:00:00', 3, 0, NOW(), NOW(), 'system'),
    ('晚班', '16:00:00', '01:00:00', 2, 0, NOW(), NOW(), 'system'),
    ('周末值班', '10:00:00', '19:00:00', 2, 1, NOW(), NOW(), 'system');

    -- 2. 初始化客服人员 (t_customer_service)
    DELETE FROM t_customer_service;
    INSERT INTO t_customer_service (cs_id, cs_name, dept_id, team_id, skill_tags, status, current_status, create_time, update_time) VALUES
    ('CS1001', '张伟', 'DEPT01', 'TEAM01', '普通咨询,订单查询', 1, 0, NOW(), NOW()),
    ('CS1002', '李娜', 'DEPT01', 'TEAM01', '售后处理,退款审核', 1, 0, NOW(), NOW()),
    ('CS1003', '王强', 'DEPT01', 'TEAM02', '投诉专员,纠纷处理', 1, 0, NOW(), NOW()),
    ('CS1004', '刘洋', 'DEPT01', 'TEAM02', 'VIP服务,大客户对接', 1, 0, NOW(), NOW()),
    ('CS1005', '陈静', 'DEPT02', 'TEAM01', '技术支持,系统故障', 1, 0, NOW(), NOW()),
    ('CS1006', '赵军', 'DEPT02', 'TEAM01', '普通咨询,账号问题', 1, 0, NOW(), NOW()),
    ('CS1007', '孙丽', 'DEPT02', 'TEAM02', '售后处理,物流追踪', 1, 0, NOW(), NOW()),
    ('CS1008', '周杰', 'DEPT02', 'TEAM02', '投诉专员,服务评价', 1, 0, NOW(), NOW()),
    ('CS1009', '吴刚', 'DEPT03', 'TEAM01', '技术支持,支付问题', 1, 0, NOW(), NOW()),
    ('CS1010', '郑敏', 'DEPT03', 'TEAM01', 'VIP服务,活动咨询', 1, 0, NOW(), NOW());

    -- 3. 初始化排班记录 (t_schedule) - 生成未来7天的排班
    DELETE FROM t_schedule;
    -- 这里仅插入示例数据，更复杂的逻辑由自动排班接口实现
    INSERT INTO t_schedule (cs_id, shift_id, schedule_date, status, create_time, update_time)
    SELECT cs_id, 
           (SELECT shift_id FROM t_shift_config LIMIT 1), 
           DATE_ADD(CURDATE(), INTERVAL 1 DAY), 
           0, NOW(), NOW()
    FROM t_customer_service LIMIT 5;

    -- 4. 初始化会话记录 (t_conversation)
    DELETE FROM t_conversation;
    INSERT INTO t_conversation (conv_id, user_id, user_nickname, cs_id, source, start_time, end_time, msg_type, is_manual_adjust, category_id, tags, is_core, status, create_time, update_time) VALUES
    (UUID(), 'USER001', '快乐小狗', 'CS1001', 'app', DATE_SUB(NOW(), INTERVAL 2 HOUR), DATE_SUB(NOW(), INTERVAL 1 HOUR), 1, 0, 1, '咨询,订单', 0, 2, NOW(), NOW()),
    (UUID(), 'USER002', '旅行达人', 'CS1002', 'web', DATE_SUB(NOW(), INTERVAL 5 HOUR), DATE_SUB(NOW(), INTERVAL 4 HOUR), 1, 0, 2, '售后,退款', 1, 2, NOW(), NOW()),
    (UUID(), 'USER003', '美食家', 'CS1001', 'wechat', DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_SUB(NOW(), INTERVAL 23 HOUR), 1, 0, 1, '咨询', 0, 2, NOW(), NOW()),
    (UUID(), 'USER004', '科技控', 'CS1005', 'app', DATE_SUB(NOW(), INTERVAL 30 MINUTE), NULL, 1, 0, 3, '技术', 0, 1, NOW(), NOW()),
    (UUID(), 'USER005', '购物狂', 'CS1003', 'web', DATE_SUB(NOW(), INTERVAL 10 MINUTE), NULL, 1, 0, 4, '投诉', 1, 1, NOW(), NOW()),
    (UUID(), 'USER006', '潜水员', 'CS1004', 'app', DATE_SUB(NOW(), INTERVAL 3 DAY), DATE_SUB(NOW(), INTERVAL 2 DAY), 1, 0, 1, 'VIP', 0, 2, NOW(), NOW()),
    (UUID(), 'USER007', '路人甲', 'CS1002', 'wechat', NOW(), NULL, 1, 0, 2, '售后', 0, 0, NOW(), NOW()),
    (UUID(), 'USER008', '测试号', 'CS1006', 'web', DATE_SUB(NOW(), INTERVAL 4 HOUR), DATE_SUB(NOW(), INTERVAL 3 HOUR), 1, 0, 1, '咨询', 0, 2, NOW(), NOW()),
    (UUID(), 'USER009', '匿名用户', 'CS1007', 'app', DATE_SUB(NOW(), INTERVAL 12 HOUR), DATE_SUB(NOW(), INTERVAL 11 HOUR), 1, 0, 2, '物流', 0, 2, NOW(), NOW()),
    (UUID(), 'USER010', '管理员', 'CS1008', 'internal', DATE_SUB(NOW(), INTERVAL 1 WEEK), DATE_SUB(NOW(), INTERVAL 6 DAY), 1, 0, 4, '投诉', 1, 2, NOW(), NOW());

END //

DELIMITER ;
