module github.com/jan-zabloudil/release-manager

go 1.22.0

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/nedpals/supabase-go v0.4.0
	go.strv.io/env v0.1.0
)

require (
	github.com/go-chi/chi/v5 v5.0.12 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/nedpals/postgrest-go v0.1.3 // indirect
)

replace github.com/nedpals/supabase-go => github.com/jan-zabloudil/supabase-go v0.0.0-20240228061618-336148d93bcd
