package utils

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Key      string // Private Key Content
}

func (c *SSHConfig) RunCommand(cmd string) (string, error) {
	var authMethods []ssh.AuthMethod

	if c.Key != "" {
		signer, err := ssh.ParsePrivateKey([]byte(c.Key))
		if err != nil {
			return "", fmt.Errorf("failed to parse private key: %v", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if c.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.Password))
	}

	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For simplicity, in production should verify host key
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("failed to run command: %v, output: %s", err, output)
	}

	return string(output), nil
}
