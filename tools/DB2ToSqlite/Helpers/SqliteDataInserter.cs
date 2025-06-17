using System.Text.Json;
using DBCD;
using DBDefsLib;
using Microsoft.Data.Sqlite;

namespace DB2ToSqliteTool.Helpers;

public static class SqliteDataInserter
{
	public static void InsertRows(IDBCDStorage storage, string tableName, Structs.DBDefinition dbDef,
		SqliteConnection connection, uint build)
	{
		var versionDef = dbDef.versionDefinitions.LastOrDefault(x => x.builds.Any(y => y.build == build));

		var columnNames = versionDef.definitions.Select(def => def.name).ToList();


		var pkDefinition = versionDef.definitions.FirstOrDefault(def => def.isID);

		var pkColumn = pkDefinition.name;


		var columnsPart = string.Join(", ", columnNames.Select(c => $"[{c}]"));
		var valuesPart = string.Join(", ", columnNames.Select(c => "@" + c));


		var updateColumns = columnNames.Where(c => c != pkColumn).ToList();
		var updateClause = updateColumns.Any()
			? "DO UPDATE SET " + string.Join(", ", updateColumns.Select(c => $"[{c}] = excluded.[{c}]"))
			: "DO NOTHING";


		var upsertSql =
			$"INSERT INTO [{tableName}] ({columnsPart}) VALUES ({valuesPart}) ON CONFLICT([{pkColumn}]) {updateClause};";


		using var transaction = connection.BeginTransaction();
		using (var command = connection.CreateCommand())
		{
			command.Transaction = transaction;

			// create indexes
			foreach (var col in columnNames)
			{
				var colDef = versionDef.definitions.FirstOrDefault(d => d.name == col);
				if (colDef.isRelation)
				{
					command.CommandText = $"CREATE INDEX IF NOT EXISTS idx_{col.ToLower()} ON {tableName} ({col});";
					try
					{
						command.ExecuteNonQuery();
					}
					catch (SqliteException se)
					{
						Console.WriteLine("Error executing command:");
						Console.WriteLine(command.CommandText);
						foreach (SqliteParameter param in command.Parameters)
							Console.WriteLine($"{param.ParameterName} = {param.Value}");
						Console.WriteLine("Exception: " + se.Message);
						throw;
					}
				}
			}


			command.CommandText = upsertSql;

			// store a list of all relations to create indexes for those columns for faster lookups
			var indexList = new HashSet<string>();
			foreach (var kvp in storage.Values)
			{
				command.Parameters.Clear();

				foreach (var col in columnNames)
				{
					var value = kvp[col] ?? DBNull.Value;
					var colDef = versionDef.definitions.FirstOrDefault(d => d.name == col);
					if (colDef.isRelation && value == (object)0) value = DBNull.Value;
					if (colDef.arrLength > 0)
						if (value is Array arr)
							value = JsonSerializer.Serialize(arr);

					command.Parameters.AddWithValue("@" + col, value);
				}

				try
				{
					command.ExecuteNonQuery();
				}
				catch (SqliteException se)
				{
					Console.WriteLine("Error executing command:");
					Console.WriteLine(command.CommandText);
					foreach (SqliteParameter param in command.Parameters)
						Console.WriteLine($"{param.ParameterName} = {param.Value}");
					Console.WriteLine("Exception: " + se.Message);
					throw;
				}
			}
		}

		transaction.Commit();
	}
}
