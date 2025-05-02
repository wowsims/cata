using DBDefsLib;
using Microsoft.Data.Sqlite;

namespace DB2ToSqliteTool.Helpers;

public static class SqliteDbCreator
{
	public static void CreateDatabaseWithDefinitions(Dictionary<string, Structs.DBDefinition> definitions,
		string sqliteFilePath, uint build)
	{
		if (File.Exists(sqliteFilePath)) File.Delete(sqliteFilePath);

		var connectionString = new SqliteConnectionStringBuilder
		{
			DataSource = sqliteFilePath,
			Mode = SqliteOpenMode.ReadWriteCreate
		}.ToString();

		using var connection = new SqliteConnection(connectionString);
		connection.Open();

		using (var pragmaCommand = connection.CreateCommand())
		{
			pragmaCommand.CommandText = "PRAGMA foreign_keys = ON;";
			pragmaCommand.ExecuteNonQuery();
		}

		using (var transaction = connection.BeginTransaction())
		{
			foreach (var kvp in definitions)
			{
				var tableName = kvp.Key;
				var dbDef = kvp.Value;

				var versionDef = dbDef.versionDefinitions.LastOrDefault(x => x.builds.Any(y => y.build == build));

				var columnDefinitionsSql = new List<string>();
				var foreignKeysSql = new List<string>();

				foreach (var def in versionDef.definitions)
				{
					if (!dbDef.columnDefinitions.TryGetValue(def.name, out var colDef))
						throw new Exception($"Column definition for {def.name} not found in table {tableName}");

					if (def.arrLength == 0)
					{
						var sqliteType = MapToSqLiteType(colDef.type);
						var nullability = !string.IsNullOrEmpty(colDef.foreignTable) &&
						                  !string.IsNullOrEmpty(colDef.foreignColumn) &&
						                  !def.isID
							? " NULL"
							: "";
						var columnSql = $"[{def.name}] {sqliteType}{nullability}";
						if (def.isID) columnSql += " PRIMARY KEY";
						columnDefinitionsSql.Add(columnSql);

						if (!string.IsNullOrEmpty(colDef.foreignTable) &&
						    !string.IsNullOrEmpty(colDef.foreignColumn))
							foreignKeysSql.Add(
								$"CREATE INDEX IF NOT EXISTS IX_{tableName}_{def.name} ON [{tableName}] ([{def.name}])");
					}
					else
					{
						var mainColumnSql = $"[{def.name}] TEXT";

						columnDefinitionsSql.Add(mainColumnSql);

						// (For example, an "int" or "uint" becomes INTEGER, "float" becomes REAL, etc.)
						var elementType = MapToSqLiteType(colDef.type);

						// Create a generated column for each array index.
						for (var i = 0; i < def.arrLength; i++)
						{
							// Generated column syntax:
							// [ColumnName_i] <Type> GENERATED ALWAYS AS (json_extract([ColumnName], '$[i]')) VIRTUAL
							var genColumn =
								$"[{def.name}_{i}] {elementType} GENERATED ALWAYS AS (json_extract([{def.name}], '$[{i}]')) VIRTUAL";
							columnDefinitionsSql.Add(genColumn);
						}
					}
				}

				var allColumns = new List<string>(columnDefinitionsSql);


				var createTableSql = $"CREATE TABLE IF NOT EXISTS [{tableName}] ({string.Join(", ", allColumns)});";

				using (var command = connection.CreateCommand())
				{
					command.CommandText = createTableSql;
					command.Transaction = transaction;
					command.ExecuteNonQuery();
				}

				foreach (var indexSql in foreignKeysSql)
				{
					using var cmd = connection.CreateCommand();
					cmd.CommandText = indexSql;
					cmd.Transaction = transaction;
					cmd.ExecuteNonQuery();
				}
			}

			transaction.Commit();
		}

		connection.Close();
	}

	private static string MapToSqLiteType(string type)
	{
		return type switch
		{
			"int" or "uint" => "INTEGER",
			"float" => "REAL",
			"string" or "locstring" => "TEXT",
			_ => throw new Exception("Unsupported type: " + type)
		};
	}
}
