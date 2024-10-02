package internal

import (
	"PortalCRG/internal/repository/entity"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/base64" // Importar paquete para decodificación Base64
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql" // Importar el driver de MySQL
)

type DriveService struct {
	db *sql.DB
}

func registerTLSConfig() error {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("ca.pem")
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo ca.pem: %v", err)
	}

	// Agregar el CA al root pool
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("falló al agregar ca.pem al cert pool")
	}

	// Crear la configuración TLS
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs: rootCertPool,
	})

	return nil
}

func parseMySQLConnectionString(connStr string) (dbUser, dbPassword, dbHost, dbPort, dbName string, err error) {
	// Verificar si el string comienza con el prefijo correcto
	if strings.HasPrefix(connStr, "mysql://") {
		connStr = strings.TrimPrefix(connStr, "mysql://")
	} else {
		return "", "", "", "", "", fmt.Errorf("el string de conexión no es válido")
	}

	// Extraer la parte de conexión (antes del '?')
	connStrParts := strings.Split(connStr, "?")
	if len(connStrParts) == 0 {
		return "", "", "", "", "", fmt.Errorf("string de conexión no válido")
	}
	connStr = connStrParts[0]

	// Parsear la URL
	parsedURL, err := url.Parse("//" + connStr) // Añadimos // para hacerla una URL válida
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error al parsear el string de conexión: %v", err)
	}

	// Obtener el nombre de usuario y la contraseña
	if parsedURL.User != nil {
		dbUser = parsedURL.User.Username()
		dbPassword, _ = parsedURL.User.Password()
	}

	// Obtener el host y puerto
	hostPort := strings.Split(parsedURL.Host, ":")
	if len(hostPort) == 2 {
		dbHost = hostPort[0]
		dbPort = hostPort[1]
	} else {
		return "", "", "", "", "", fmt.Errorf("el host y puerto no son válidos en el string de conexión")
	}

	// Obtener el nombre de la base de datos
	dbName = strings.TrimPrefix(parsedURL.Path, "/")

	return dbUser, dbPassword, dbHost, dbPort, dbName, nil
}

// Conectar a la base de datos usando las variables de entorno
func (d *DriveService) Connect() error {
	// Registrar la configuración TLS antes de crear la conexión
	err := registerTLSConfig()
	if err != nil {
		return err
	}

	connStr := os.Getenv("MYSQLROMS_CONNECTION_STRING")

	dbUser, dbPassword, dbHost, dbPort, dbName, err := parseMySQLConnectionString(connStr)
	if err != nil {
		log.Println("Error:", err)
		return err
	}

	fmt.Printf("dbUser: %s\n", dbUser)
	fmt.Printf("dbHost: %s\n", dbHost)
	fmt.Printf("dbPort: %s\n", dbPort)
	fmt.Printf("dbName: %s\n", dbName)

	// Crear el string de conexión, usando `tls=custom` para nuestra configuración TLS
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=custom", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	// Verificar que la conexión funcione
	err = db.Ping()
	if err != nil {
		return err
	}

	d.db = db

	return nil
}

// Crear la tabla retro_files con id_tips como clave primaria y content como BLOB
func (d *DriveService) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS retro_files (
		id_tips VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		size INT NOT NULL,
		content LONGBLOB  NOT NULL
	);`
	_, err := d.db.Exec(query)
	return err
}

// Crear un nuevo RetroFile
func (d *DriveService) CreateFile(retroFile *entity.RetroFile) error {
	// Decodificar el contenido Base64 a bytes
	decodedContent, err := base64.StdEncoding.DecodeString(retroFile.Content)
	if err != nil {
		return fmt.Errorf("error al decodificar el contenido: %v", err)
	}

	query := `INSERT INTO retro_files (id_tips, name, type, size, content) VALUES (?, ?, ?, ?, ?)`
	_, err = d.db.Exec(query, retroFile.IdTips, retroFile.Name, retroFile.Type, retroFile.Size, decodedContent)
	if err != nil {
		return err
	}
	return nil
}

// Leer un RetroFile por id_tips
func (d *DriveService) GetFileByID(idTips string) (*entity.RetroFile, error) {
	query := `SELECT name, type, size, content FROM retro_files WHERE id_tips = ?`
	row := d.db.QueryRow(query, idTips)

	var retroFile entity.RetroFile
	retroFile.IdTips = idTips // Asignamos el id_tips que estamos buscando

	if err := row.Scan(&retroFile.Name, &retroFile.Type, &retroFile.Size, &retroFile.Content); err != nil {
		return nil, err
	}

	// Convertir el contenido a Base64 antes de devolverlo
	//retroFile.Content = base64.StdEncoding.EncodeToString([]byte(retroFile.Content))

	return &retroFile, nil
}

// Actualizar un RetroFile
func (d *DriveService) UpdateFile(retroFile *entity.RetroFile) error {
	// Decodificar el contenido Base64 a bytes
	decodedContent, err := base64.StdEncoding.DecodeString(retroFile.Content)
	if err != nil {
		return fmt.Errorf("error al decodificar el contenido: %v", err)
	}

	query := `UPDATE retro_files SET name = ?, type = ?, size = ?, content = ? WHERE id_tips = ?`
	_, err = d.db.Exec(query, retroFile.Name, retroFile.Type, retroFile.Size, decodedContent, retroFile.IdTips)
	return err
}

// Eliminar un RetroFile por id_tips
func (d *DriveService) DeleteFile(idTips string) error {
	query := `DELETE FROM retro_files WHERE id_tips = ?`
	_, err := d.db.Exec(query, idTips)
	return err
}

// SaveFile realiza un INSERT si no existe el archivo, o un UPDATE si ya existe
func (d *DriveService) SaveFile(retroFile *entity.RetroFile) error {
	// Verificar si el archivo ya existe en la base de datos
	queryCheck := `SELECT COUNT(*) FROM retro_files WHERE id_tips = ?`
	var count int
	err := d.db.QueryRow(queryCheck, retroFile.IdTips).Scan(&count)
	if err != nil {
		return fmt.Errorf("error al verificar existencia de archivo: %v", err)
	}

	if count == 0 {
		// Si no existe, usar el método CreateFile
		err = d.CreateFile(retroFile)
		if err != nil {
			return fmt.Errorf("error al insertar archivo: %v", err)
		}
		log.Println("Archivo insertado correctamente")
	} else {
		// Si existe, usar el método UpdateFile
		err = d.UpdateFile(retroFile)
		if err != nil {
			return fmt.Errorf("error al actualizar archivo: %v", err)
		}
		log.Println("Archivo actualizado correctamente")
	}

	return nil
}
