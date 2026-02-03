package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"api_kino/config/app"
	"api_kino/config/database"
)

func GenerateModel() {
	dbConfig := database.DBConfig()
	dbHost := dbConfig.Host
	dbPort := dbConfig.Port
	dbDatabase := dbConfig.Database
	dbUsername := dbConfig.Username
	dbPassword := dbConfig.Password
	logType := logger.Error
	if app.Config().GinMode == gin.DebugMode {
		logType = logger.Info
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./app/models/org/query",
		ModelPkgPath:  "./app/models/org/model",
		FieldNullable: true,
		/* Mode: gen.WithoutContext|gen.WithDefaultQuery*/
		//if you want the nullable field generation property to be pointer type, set FieldNullable true
		/* FieldNullable: true,*/
		//if you want to assign field which has default value in `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		/* FieldCoverable: true,*/
		// if you want generate field with unsigned integer type, set FieldSignable true
		/* FieldSignable: true,*/
		//if you want to generate index tags from database, set FieldWithIndexTag true
		/* FieldWithIndexTag: true,*/
		//if you want to generate type tags from database, set FieldWithTypeTag true
		/* FieldWithTypeTag: true,*/
		//if you need unit tests for query code, set WithUnitTest true
		/* WithUnitTest: true, */
	})
	// dataMap mapping relationship
	var dataMap = map[string]func(detailType string) (dataType string){
		"json": func(detailType string) (dataType string) { return "datatypes.JSON" },
	}
	gormConfig := &gorm.Config{
		Logger:                 logger.Default.LogMode(logType),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d application_name=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUsername, dbPassword, dbDatabase, dbPort, "api_kino")
	connection, _ := gorm.Open(postgres.New(postgres.Config{
		DSN: connString,
	}), gormConfig)
	g.UseDB(connection)
	g.WithDataTypeMap(dataMap)
	g.ApplyBasic(g.GenerateAllTable()...)

	userTb := g.GenerateModel("users")
	userSessionTb := g.GenerateModel("users_session")
	unitTypeTb := g.GenerateModel("unit_type")
	unitTb := g.GenerateModel("unit")
	positionTb := g.GenerateModel("position")
	positionTypeTb := g.GenerateModel("position_type")
	draftTransferPositionTb := g.GenerateModel("draft_transfer_position")
	draftTransferPositionPersonTb := g.GenerateModel("draft_transfer_position_person")
	draftTransferRankTb := g.GenerateModel("draft_transfer_rank")
	draftTransferRankPersonTb := g.GenerateModel("draft_transfer_rank_person")
	positionMapTb := g.GenerateModel("position_map",
		gen.FieldIgnore("c1"),
		gen.FieldIgnore("c2"),
		gen.FieldIgnore("c3"),
		gen.FieldIgnore("c4"),
		gen.FieldIgnore("c5"),
		gen.FieldIgnore("c6"),
		gen.FieldIgnore("c7"),
		gen.FieldIgnore("c8"),
		gen.FieldIgnore("c9"),
		gen.FieldIgnore("c10"),
		gen.FieldIgnore("c11"),
		gen.FieldIgnore("c12"),
		gen.FieldIgnore("c13"),
		gen.FieldIgnore("c14"),
	)
	unitMapTb := g.GenerateModel("unit_map") //gen.FieldIgnore("c1"),
	//gen.FieldIgnore("c2"),

	unit := g.GenerateModel("unit",
		gen.FieldRelate(field.BelongsTo, "RefMapData", unitMapTb,
			&field.RelateConfig{
				RelatePointer: true, //nullable
				GORMTag:       "foreignKey:code;references:ref_map",
			}),
		gen.FieldRelate(field.BelongsTo, "UnitType", unitTypeTb,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //nullable
				GORMTag:       "foreignKey:id;references:unit_type_id",
				//GORMTag: "foreignKey:unit_type_id",
			}),
		gen.FieldRelate(field.BelongsTo, "Unit", unitTb,
			&field.RelateConfig{
				//RelateSlicePointer:true,
				RelatePointer: true, //nullable
				//RelateSlice: true,
				GORMTag: "foreignKey:id;references:unit_id",
				//GORMTag: "foreignKey:unit_id",
			}),
		gen.FieldRelate(field.HasMany, "Positions", positionTb,
			&field.RelateConfig{
				RelateSlice: true,
				//RelatePointer: true, //nullable
				//RelateSlice: true,
				GORMTag: "foreignKey:unit_id;references:id",
			}),
	)
	user := g.GenerateModel("users",
		gen.FieldJSONTag("password", "-"),
		gen.FieldRelate(field.HasMany, "UserSession", userSessionTb,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag:     "foreignKey:users_id;references:id",
			}),
		gen.FieldRelate(field.BelongsTo, "Unit", unit,
			&field.RelateConfig{
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:unit_id",
			}),
	)
	userSession := g.GenerateModel("users_session",
		gen.FieldRelate(field.BelongsTo, "User", userTb,
			&field.RelateConfig{
				RelatePointer: true, //not nullable
				//GORMTag: "foreignKey:users_id",
				GORMTag: "foreignKey:id;references:users_id",
			}),
	)
	position := g.GenerateModel("position",
		gen.FieldRelate(field.BelongsTo, "RefMapData", positionMapTb,
			&field.RelateConfig{
				RelatePointer: true, //nullable
				GORMTag:       "foreignKey:code;references:ref_map",
			}),
		gen.FieldRelate(field.BelongsTo, "Unit", unit,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:unit_id",
			}),
		gen.FieldRelate(field.BelongsTo, "PositionType", positionTypeTb,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:position_type_id",
			}),
		gen.FieldRelate(field.BelongsTo, "Position", positionTb,
			&field.RelateConfig{
				RelatePointer: true, //nullable
				GORMTag:       "foreignKey:id;references:position_id",
			}),
	)
	draftTransferPosition := g.GenerateModel("draft_transfer_position",
		gen.FieldRelate(field.HasMany, "DraftTransferPositionPerson", draftTransferPositionPersonTb,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag:     "foreignKey:draft_transfer_position_id;references:id",
			}),
	)
	draftTransferPositionPerson := g.GenerateModel("draft_transfer_position_person",
		gen.FieldRelate(field.BelongsTo, "DraftTransferPosition", draftTransferPositionTb,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:draft_transfer_position_id",
			}),
	)
	draftTransferRank := g.GenerateModel("draft_transfer_rank",
		gen.FieldRelate(field.HasMany, "DraftTransferPositionPerson", draftTransferRankPersonTb,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag:     "foreignKey:draft_transfer_rank_id;references:id",
			}),
	)
	draftTransferRankPerson := g.GenerateModel("draft_transfer_rank_person",
		gen.FieldRelate(field.BelongsTo, "DraftTransferPosition", draftTransferRankTb,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:draft_transfer_rank_id",
			}),
	)
	draftProfileUpdateTb := g.GenerateModel("draft_personnel_update")
	projectTb := g.GenerateModel("project")
	projectContentTb := g.GenerateModel("project_content")
	projectContentTagTb := g.GenerateModel("project_content_tag")
	projectAttachmentTb := g.GenerateModel("project_attachment")
	project := g.GenerateModel("project",
		gen.FieldRelate(field.BelongsTo, "Project", projectTb,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:project_id",
			}),
		gen.FieldRelate(field.HasMany, "ProjectContent", projectContentTb,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag:     "foreignKey:project_id;references:id",
			}),
		gen.FieldRelate(field.HasMany, "ProjectAttachment", projectAttachmentTb,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag:     "foreignKey:project_id;references:id",
			}),
	)
	projectContent := g.GenerateModel("project_content",
		gen.FieldRelate(field.BelongsTo, "Project", projectTb,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:project_id",
			}),
		gen.FieldRelate(field.HasMany, "ProjectContentTag", projectContentTagTb,
			&field.RelateConfig{
				RelateSlice: true,
				GORMTag:     "foreignKey:project_content_id;references:id",
			}),
	)
	projectContentTag := g.GenerateModel("project_content_tag",
		gen.FieldRelate(field.BelongsTo, "ProjectContent", projectContentTb,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:project_content_id",
			}),
	)
	projectAttachment := g.GenerateModel("project_attachment",
		gen.FieldRelate(field.BelongsTo, "Project", projectTb,
			&field.RelateConfig{
				//RelateSlice: true,
				RelatePointer: true, //not nullable
				GORMTag:       "foreignKey:id;references:project_id",
			}),
	)
	g.ApplyBasic(positionMapTb, unitMapTb, unitTypeTb, positionTypeTb, user, unit, position, userSession, draftTransferPosition, draftTransferPositionPerson, draftTransferRank, draftTransferRankPerson, draftProfileUpdateTb, project, projectContent, projectContentTag, projectAttachment)
	g.Execute()
}
