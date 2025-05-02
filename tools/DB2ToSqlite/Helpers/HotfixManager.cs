using System.Text.Json;
using DBCD.IO;

//https://github.com/Marlamin/wow.tools.local/blob/main/Services/HotfixManager.cs
namespace DB2ToSqliteTool.Helpers;

public static class HotfixManager
{
	public static Dictionary<uint, HotfixReader> HotfixReaders = [];
	public static Dictionary<uint, List<DBCacheParser>> DbcacheParsers = [];
	public static Dictionary<int, DateTime>? PushIdDetected = [];

	public static Dictionary<uint, string> TableNames = Directory.EnumerateFiles("DBDCache/")
		.ToDictionary(x => Hash(Path.GetFileNameWithoutExtension(x).ToUpper()),
			x => Path.GetFileNameWithoutExtension(x));

	private static void LoadPushIDs()
	{
		PushIdDetected =
			JsonSerializer.Deserialize<Dictionary<int, DateTime>>(File.ReadAllText("knownPushIDs.json"));
	}

	private static void SavePushIDs()
	{
		File.WriteAllText("knownPushIDs.json", JsonSerializer.Serialize(PushIdDetected));
	}

	public static void LoadCaches(string wowLocation)
	{
		if (!File.Exists("knownPushIDs.json"))
			SavePushIDs();

		LoadPushIDs();

		Console.WriteLine("Reloading all hotfixes..");
		HotfixReaders.Clear();

		if (Directory.Exists("caches"))
			foreach (var file in Directory.GetFiles("caches", "*.bin", SearchOption.AllDirectories))
			{
				var reader = new HotfixReader(file);
				if (!HotfixReaders.ContainsKey((uint)reader.BuildId))
					HotfixReaders.Add((uint)reader.BuildId, reader);

				HotfixReaders[(uint)reader.BuildId].CombineCache(file);

				if (!DbcacheParsers.ContainsKey((uint)reader.BuildId))
					DbcacheParsers.Add((uint)reader.BuildId, []);

				var newCache = new DBCacheParser(file);
				DbcacheParsers[(uint)reader.BuildId].Add(newCache);

				var newPushIDs = newCache.hotfixes.Where(x => x.pushID > 0 && PushIdDetected != null && !PushIdDetected.ContainsKey(x.pushID))
					.Select(x => x.pushID).ToList();
				foreach (var newPushId in newPushIDs)
				{
					PushIdDetected?.TryAdd(newPushId, DateTime.Now);
					Console.WriteLine("Detected new pushID " + newPushId + " at " + DateTime.Now.ToShortTimeString());
				}

				Console.WriteLine("Loaded hotfixes from caches directory for build " + reader.BuildId);
			}


		foreach (var file in Directory.GetFiles(wowLocation, "DBCache.bin", SearchOption.AllDirectories))
		{
			var reader = new HotfixReader(file);
			if (!HotfixReaders.ContainsKey((uint)reader.BuildId))
				HotfixReaders.Add((uint)reader.BuildId, reader);

			HotfixReaders[(uint)reader.BuildId].CombineCache(file);

			if (!DbcacheParsers.ContainsKey((uint)reader.BuildId))
				DbcacheParsers.Add((uint)reader.BuildId, []);

			var newCache = new DBCacheParser(file);
			DbcacheParsers[(uint)reader.BuildId].Add(newCache);

			var newPushIDs = newCache.hotfixes.Where(x => x.pushID > 0 && PushIdDetected != null && !PushIdDetected.ContainsKey(x.pushID))
				.Select(x => x.pushID).ToList();
			foreach (var newPushId in newPushIDs)
			{
				PushIdDetected?.TryAdd(newPushId, DateTime.Now);
				Console.WriteLine("Detected new pushID " + newPushId + " at " + DateTime.Now.ToShortTimeString());
			}

			Console.WriteLine("Loaded hotfixes from client for build " + reader.BuildId);
		}

		SavePushIDs();
	}

	public static void Clear()
	{
		HotfixReaders.Clear();
		DbcacheParsers.Clear();
	}

	private static uint Hash(string s)
	{
		var sHashtable = new uint[]
		{
			0x486E26EE, 0xDCAA16B3, 0xE1918EEF, 0x202DAFDB,
			0x341C7DC7, 0x1C365303, 0x40EF2D37, 0x65FD5E49,
			0xD6057177, 0x904ECE93, 0x1C38024F, 0x98FD323B,
			0xE3061AE7, 0xA39B0FA1, 0x9797F25F, 0xE4444563
		};

		uint v = 0x7fed7fed;
		var x = 0xeeeeeeee;
		for (var i = 0; i < s.Length; i++)
		{
			var c = (byte)s[i];
			v += x;
			v ^= sHashtable[(c >> 4) & 0xf] - sHashtable[c & 0xf];
			x = x * 33 + v + c + 3;
		}

		return v;
	}
}
