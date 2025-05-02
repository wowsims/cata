using TACTSharp;

namespace DB2ToSqliteTool.Helpers;

public class BindableSettings : Settings
{
	public new string Region
	{
		get => base.Region;
		set => base.Region = value;
	}

	public new string Product
	{
		get => base.Product;
		set => base.Product = value;
	}

	public new RootInstance.LocaleFlags Locale
	{
		get => base.Locale;
		set => base.Locale = value;
	}

	public new RootInstance.LoadMode RootMode
	{
		get => base.RootMode;
		set => base.RootMode = value;
	}

	public new string? BaseDir
	{
		get => base.BaseDir;
		set => base.BaseDir = value;
	}

	public new string? BuildConfig
	{
		get => base.BuildConfig;
		set => base.BuildConfig = value;
	}

	public new string? CDNConfig
	{
		get => base.CDNConfig;
		set => base.CDNConfig = value;
	}

	public new string CacheDir
	{
		get => base.CacheDir;
		set => base.CacheDir = value;
	}

	public new bool ListfileFallback
	{
		get => base.ListfileFallback;
		set => base.ListfileFallback = value;
	}

	public new string ListfileURL
	{
		get => base.ListfileURL;
		set => base.ListfileURL = value;
	}

	public List<string> GameTables { get; set; } = [];
	public string GameTablesOutDirectory { get; set; } = "";
}
