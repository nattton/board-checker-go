# bkzy-organizer

Board Checker

# Build
	./build/build.sh 

# Start Web

	./bin/web

## Usage of web:

  -addr string
  
	HTTP Network Address (default ":4000")

  -dsn string
  
	Database DSN (default "$BC_DSN")
  
  -html-dir string
  
	Path to static assets (default "$GOPATH/src/gitlab.com/code-mobi/board-checker/ui/html")
    	
  -secret string
  
	Secret key
    	
  -static-dir string
  
	Path to static assets (default "$GOPATH/src/gitlab.com/code-mobi/board-checker/ui/static")