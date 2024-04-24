Run the following for *developer instance*:

```bash
go run main.go initialize_organization;
go run main.go import_tag;
go run main.go import_videocategory;
go run main.go import_exercise;
go run main.go import_equipment;
go run main.go import_bp8exercise;
go run main.go import_bp8video --prefix=dev;
go run main.go import_bp8thumbnail --prefix=dev;
```
