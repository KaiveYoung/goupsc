package nut

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type NUTClient struct {
	conn   net.Conn
	reader *bufio.Scanner
	ups    string
}

// Connect 连接到 NUT 服务器
func (c *NUTClient) Connect(host string, port string, ups string) error {
	var err error
	c.conn, err = net.Dial("tcp", host+":"+port)
	if err != nil {
		return err
	}
	c.reader = bufio.NewScanner(c.conn)
	c.ups = ups
	return nil
}

// SendCommand 发送命令到 NUT 服务器
func (c *NUTClient) SendCommand(cmd string) error {
	_, err := fmt.Fprintf(c.conn, cmd+"\n")
	return err
}

// ReadResponse 读取服务器响应
func (c *NUTClient) ReadResponse() ([]string, error) {
	var lines []string
	inResponse := false
	for c.reader.Scan() {
		line := c.reader.Text()

		// 判断是否开始读取响应
		if strings.HasPrefix(line, "BEGIN ") {
			inResponse = true
			continue
		}

		// 判断是否结束读取响应
		if strings.HasPrefix(line, "END ") {
			break
		}

		// 读取响应内容
		if inResponse {
			lines = append(lines, line)
		}
	}

	if err := c.reader.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

// Close 关闭连接
func (c *NUTClient) Close() {
	c.conn.Close()
}

func (c *NUTClient) ListVars() ([]string, error) {
	err := c.SendCommand(fmt.Sprintf("LIST VAR %s", c.ups))
	if err != nil {
		return nil, err
	}
	return c.ReadResponse()
}

type VarValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *NUTClient) ParseVarValue(varValue string) (VarValue, error) {
	v := VarValue{}
	perfix := fmt.Sprintf("VAR %s ", c.ups)
	if !strings.HasPrefix(varValue, perfix) {
		return v, fmt.Errorf("line does not start with %s: %s", perfix, varValue)
	}
	varValue = strings.TrimPrefix(varValue, perfix)
	quoteIndex := strings.Index(varValue, `"`)
	if quoteIndex == -1 {
		return v, fmt.Errorf("no starting quote found in line: %s", varValue)
	}
	v.Name = strings.TrimSpace(varValue[:quoteIndex])
	v.Value = strings.Trim(varValue[quoteIndex+1:], `"`)
	return v, fmt.Errorf("数据错误")
}
