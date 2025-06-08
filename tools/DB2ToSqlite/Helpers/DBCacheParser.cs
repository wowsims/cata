﻿namespace DB2ToSqliteTool.Helpers;

public struct HotfixEntry
{
	public uint regionID;
	public int pushID;
	public uint uniqueID;
	public uint tableHash;
	public uint recordID;
	public int dataSize;
	public byte status;
	public byte[] data;
}

//https://github.com/Marlamin/wow.tools.local/blob/main/Services/DBCacheParser.cs
public class DBCacheParser
{
	public int build;
	public List<HotfixEntry> hotfixes = [];

	public DBCacheParser(string filename)
	{
		using (var fs = File.Open(filename, FileMode.Open, FileAccess.Read, FileShare.ReadWrite))
		using (var bin = new BinaryReader(fs))
		{
			var hotfix = new HotfixEntry();
			bin.ReadUInt32(); // Signature
			var version = bin.ReadUInt32();
			if (version != 9)
				//Console.WriteLine("Unsupported DBCache version " + version + ", skipping");
				return;
			build = bin.ReadInt32();
			bin.BaseStream.Position += 32;

			while (bin.BaseStream.Position < bin.BaseStream.Length)
			{
				bin.ReadUInt32(); // Signature
				hotfix.regionID = bin.ReadUInt32();
				hotfix.pushID = bin.ReadInt32();
				hotfix.uniqueID = bin.ReadUInt32();
				hotfix.tableHash = bin.ReadUInt32();
				hotfix.recordID = bin.ReadUInt32();
				hotfix.dataSize = bin.ReadInt32();
				hotfix.status = bin.ReadByte();
				bin.ReadBytes(3);
				hotfix.data = bin.ReadBytes(hotfix.dataSize);
				hotfixes.Add(hotfix);
			}
		}
	}
}
