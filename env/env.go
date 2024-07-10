package env

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type EnvManager struct {
	dbHost 		string;
	dbUser		string;
	dbName		string;
	dbPassword	string;
	sslmode		string;
	DbSource	string;
	DbDriver 	string;
}

func NewEnvManager() *EnvManager {
	return &EnvManager{}
}

func (env *EnvManager)dataSourceName() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
	 env.dbHost, env.dbUser, env.dbPassword, env.dbName, env.sslmode)
}

func SetupEnv(file string) *EnvManager {
	env := NewEnvManager()
	env.loadEnv(file)
	env.DbSource = env.dataSourceName()
	return env
}

func (e *EnvManager) loadEnv(file string) {
	
	fd, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error opening file, file not found")
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	
	m := make(map[string]string)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "=")
		if len(line) != 2 {
			fmt.Printf("Invalid line: %s\n", scanner.Text())
			continue
		}
		m[strings.TrimSpace(line[0])] = strings.TrimSpace(line[1])
	}
	
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}


	e.dbHost = m["DBHOST"]
	e.dbUser = m["DBUSER"]
	e.dbName = m["DBNAME"]
	e.dbPassword = m["DBPASS"]
	e.sslmode = m["SSLMODE"]
	e.DbDriver = m["DBDRIVER"]
}
