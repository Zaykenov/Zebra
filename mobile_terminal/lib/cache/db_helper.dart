// import 'dart:async';
// import 'dart:io';

// import 'package:path/path.dart' as p;
// import 'package:path_provider/path_provider.dart';
// import 'package:sqflite/sqflite.dart';

// class CacheItem {
//   String flUrl;
//   String flContentType;
//   String flContentInBase64;
//   DateTime flTimeStamp;

//   CacheItem(
//       this.flUrl, this.flContentType, this.flContentInBase64, this.flTimeStamp);
// }

// class DbHelper {
//   Database db;
//   DbHelper(this.db);

//   static DbHelper? _instance;
//   static Future<DbHelper> getInstance() async {
//     if (_instance == null) {
//       var db = await initDB();
//       _instance = DbHelper(db);
//     }
//     return _instance!;
//   }

//   static Future<Database> initDB() async {
//     Directory documentsDirectory = await getApplicationDocumentsDirectory();
//     String path = p.join(documentsDirectory.path, "zebraCrm.db");

//     return await openDatabase(path, version: 1, onOpen: (db) {},
//         onCreate: (Database db, int version) async {
//       await db.execute("CREATE TABLE TbCache ("
//           "flUrl TEXT PRIMARY KEY,"
//           "flContentType TEXT,"
//           "flContentInBase64 TEXT,"
//           "flTimeStamp TEXT"
//           ")");
//     });
//   }

//   Future<void> addCache(CacheItem cacheItem) async {
//     await db.transaction((txn) async {
//       txn.delete("TbCache", where: 'flurl = ?', whereArgs: [cacheItem.flUrl]);

//       await txn.rawInsert(
//           "INSERT Into TbCache (flUrl, flContentType, flContentInBase64, flTimeStamp)"
//           " VALUES (?,?,?,?)",
//           [
//             cacheItem.flUrl,
//             cacheItem.flContentType,
//             cacheItem.flContentInBase64,
//             DateTime.now().toIso8601String()
//           ]);
//     });
//   }

//   Future<void> clearAllCache() async {
//     await db.transaction((txn) async {
//       txn.delete("TbCache");
//     });
//   }

//   Future<CacheItem?> getFromCacheOrNull(String url,
//       {int cacheToleranseInMinutes = 8 * 60}) async {
//     var sqlJsonRows =
//         await db.query("TbCache", where: 'flurl = ?', whereArgs: [url]);
//     if (sqlJsonRows.isEmpty) {
//       return null;
//     }

//     var row = sqlJsonRows[0];
//     var cacheDate = DateTime.parse(row["flTimeStamp"].toString());

//     final differenceInMinuts = DateTime.now().difference(cacheDate).inMinutes;
//     if (differenceInMinuts > cacheToleranseInMinutes) {
//       return null;
//     }

//     return CacheItem(row["flUrl"].toString(), row["flContentType"].toString(),
//         row["flContentInBase64"].toString(), cacheDate);
//   }

//   Future<int> getCacheKeysCount() async {
//     var sqlJsonRows = await db.query("TbCache");
//     return sqlJsonRows.length;
//   }
// }
