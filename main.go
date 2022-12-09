package main

import (
	"context"
	"fmt"
	"os"

	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

var (
	pId = "project_id"
	oId = "operation_id"
)

func main() {
	ctx := context.Background()

	status, err := GetOperationStatus(ctx, pId, oId)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(status)
}

func ListInstances(ctx context.Context, projectId string) ([]*sqladmin.DatabaseInstance, error) {
	// Create the Google Cloud SQL service.
	service, err := sqladmin.NewService(ctx)
	if err != nil {
		return nil, err
	}

	// List instances for the project ID.
	instances, err := service.Instances.List(projectId).Do()
	if err != nil {
		return nil, err
	}
	return instances.Items, nil
}

func RequestExport(ctx context.Context, projectId string, instanceId string, fileName string) error {
	// Create the Google Cloud SQL service.
	service, err := sqladmin.NewService(ctx)
	if err != nil {
		return err
	}

	exportRequest := &sqladmin.InstancesExportRequest{
		ExportContext: &sqladmin.ExportContext{
			Uri:              "google_cloud_storage_uri",
			Databases:        []string{"google_cloud_sql_database_name"},
			Kind:             "sql#exportContext",
			CsvExportOptions: &sqladmin.ExportContextCsvExportOptions{SelectQuery: "select * from tokens;"},
			FileType:         "CSV",
			Offload:          true,
		},
	}

	operation, err := service.Instances.Export(projectId, instanceId, exportRequest).Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(operation.Name)
	fmt.Println(operation.Status)

	return nil
}

func GetOperationStatus(ctx context.Context, projectId string, operationName string) (string, error) {
	service, err := sqladmin.NewService(ctx)
	if err != nil {
		return "", err
	}

	operation, err := service.Operations.Get(projectId, operationName).Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// クエリ実行中にエラーが生じたなら、OperationErrors型の変数が入る。そうでないなら、nilになる。
	fmt.Println(operation.Error.Errors)
	for _, err := range operation.Error.Errors {
		fmt.Println(err)
	}

	return operation.Status, nil
}
