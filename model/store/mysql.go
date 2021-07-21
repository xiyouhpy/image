package store

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// MysqlInterface 接口整理
type MysqlInterface interface {
	// NewMysql 获取 mysql 对象
	NewMysql(host string, port string, user string, passwd string, dbname string) (*MysqlCli, error)

	// Select mysql select 方法，select 操作比较特殊，需要拿到拼接完整的 sql 作为入参
	Select(strSql string) ([]interface{}, error)
	// Insert mysql insert 操作
	Insert(strSql string, args ...interface{}) (int64, error)
	// Update mysql update 操作
	Update(strSql string, args ...interface{}) (int64, error)
	// Delete mysql delete 操作
	Delete(strSql string, args ...interface{}) (int64, error)
}

// MysqlCli mysql 对象结构
type MysqlCli struct {
	client *sql.DB
}

// MysqlConfig mysql 配置文件结构
type MysqlConfig struct {
	host   string
	port   string
	user   string
	passwd string
	dbname string
}

// NewMysql 获取 mysql 对象并初始化，使用 utf-8 编码
func (mc *MysqlConfig) NewMysql() (*MysqlCli, error) {
	if mc.host == "" || mc.port == "" || mc.user == "" || mc.passwd == "" || mc.dbname == "" {
		logrus.Warnf("NewMysql params err")
		return nil, errors.New("params err")
	}

	// mysql 服务连接
	dbServer := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", mc.user, mc.passwd, mc.host, mc.port, mc.dbname)
	client, err := sql.Open("mysql", dbServer)
	if err != nil {
		logrus.Warnf("mysql Open err, err:%s", err.Error())
		return nil, err
	}

	if err = client.Ping(); err != nil {
		logrus.Warnf("ping mysql err, err:%s", err.Error())
		return nil, err
	}
	logrus.Infof("connect to mysql, %s:%s", mc.host, mc.port)

	return &MysqlCli{client: client}, nil
}

// Select mysql select 操作，select 操作比较特殊，需要拿到拼接完整的 sql 作为入参
func (conn *MysqlCli) Select(strSql string) ([]interface{}, error) {
	if strSql == "" {
		logrus.Warnf("Select params sql err")
		return nil, errors.New("params sql err")
	}

	// 执行查询 sql 命令
	rows, err := conn.client.Query(strSql)
	if err != nil {
		logrus.Warnf("Select query err, sql:%s", strSql)
		return nil, err
	}
	defer rows.Close()

	// 保存查询结果信息
	var arrResult []interface{}

	// 循环遍历返回结果
	for rows.Next() {
		// 获取记录字段，把字段参数值和字段地址关联
		columns, _ := rows.Columns()
		arrScanArgs := make([]interface{}, len(columns))
		arrScanValue := make([]interface{}, len(columns))
		for i := range arrScanValue {
			arrScanArgs[i] = &arrScanValue[i]
		}

		// 将数据保存到 record 字典
		err = rows.Scan(arrScanArgs...)
		if err != nil {
			logrus.Warnf("Select Scan err, sql:%s", strSql)
			return nil, err
		}
		arrTemp := make(map[string]interface{})
		for i, col := range arrScanValue {
			if col != nil {
				arrTemp[columns[i]] = string(col.([]byte))
			}
		}
		arrResult = append(arrResult, arrTemp)
	}

	return arrResult, nil
}

// Insert mysql insert 操作
func (conn *MysqlCli) Insert(strSql string, args ...interface{}) (int64, error) {
	if strSql == "" {
		logrus.Warnf("Insert params sql err")
		return 0, errors.New("params sql err")
	}

	// 执行查询 sql 命令
	stmt, err := conn.client.Prepare(strSql)
	if err != nil {
		logrus.Warnf("Insert Prepare err, sql:%s", strSql)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		logrus.Warnf("Insert Exec err, sql:%s", strSql)
		return 0, err
	}

	insertId, err := res.LastInsertId()
	if err != nil {
		logrus.Warnf("Insert LastInsertId err, sql:%s", strSql)
		return 0, err
	}

	return insertId, nil
}

// Update mysql update 操作
func (conn *MysqlCli) Update(strSql string, args ...interface{}) (int64, error) {
	if strSql == "" {
		logrus.Warnf("Update params sql err")
		return 0, errors.New("params sql err")
	}

	// 执行查询 sql 命令
	stmt, err := conn.client.Prepare(strSql)
	if err != nil {
		logrus.Warnf("Update Prepare err, sql:%s", strSql)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		logrus.Warnf("Update Exec err, sql:%s", strSql)
		return 0, err
	}

	updateNum, err := res.RowsAffected()
	if err != nil {
		logrus.Warnf("Update RowsAffected err, sql:%s", strSql)
		return 0, err
	}

	return updateNum, nil
}

// Delete mysql delete 操作
func (conn *MysqlCli) Delete(strSql string, args ...interface{}) (int64, error) {
	if strSql == "" {
		logrus.Warnf("Delete params sql err")
		return 0, errors.New("params sql err")
	}

	// 执行查询 sql 命令
	stmt, err := conn.client.Prepare(strSql)
	if err != nil {
		logrus.Warnf("Delete Prepare err, sql:%s", strSql)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		logrus.Warnf("Delete Exec err, sql:%s", strSql)
		return 0, err
	}

	deleteNum, err := res.RowsAffected()
	if err != nil {
		logrus.Warnf("Delete RowsAffected err, sql:%s", strSql)
		return 0, err
	}

	return deleteNum, nil
}
