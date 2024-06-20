import 'dart:io';

import 'package:mobile_terminal/cache/enum_helper.dart';
import 'package:path_provider/path_provider.dart';
import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart';

class DbHelper {
  DbHelper._();

  static final DbHelper db = DbHelper._();

  Database? _database;

  Future<Database> get database async {
    if (_database != null) return _database!;
    _database = await initDB();
    return _database!;
  }

  initDB() async {
    Directory documentsDirectory = await getApplicationDocumentsDirectory();
    String path = join(documentsDirectory.path, "dbZebraCrm.db");

    return await openDatabase(path, version: 1, onOpen: (db) {},
        onCreate: (Database db, int version) async {
      await db.execute('''
    CREATE TABLE ${EnumHelper.enumToStr(Tables.tbSendOrdersToApiJobs)} (
      id INTEGER PRIMARY KEY,
      time INTEGER,
      request TEXT,
      response TEXT,
      exception TEXT,
      retryCount INTEGER,
      status TEXT,
      authorization TEXT,
      idempotencyKey TEXT
    )
  ''');
    });
  }

  Future insertJob(SendOrderToApiJob job) async {
    final db = await database;

    await db.transaction((txn) async {
      await txn.rawInsert(
          "INSERT Into ${EnumHelper.enumToStr(Tables.tbSendOrdersToApiJobs)} (time, request, response, exception, retryCount, status, authorization, idempotencyKey)"
          " VALUES (?,?,?,?,?,?,?,?)",
          [
            job.time.toUtc().millisecondsSinceEpoch,
            job.request,
            job.response,
            job.exception,
            job.retryCount,
            job.status,
            job.authorization,
            job.idempotencyKey
          ]);
    });
  }

  Future updateJob(
      int id, int retryCount, String status, String response) async {
    final db = await database;

    await db.rawUpdate(
        'UPDATE ${EnumHelper.enumToStr(Tables.tbSendOrdersToApiJobs)} SET retryCount = ?, status = ?, response = ? WHERE id = ?',
        [retryCount, status, response, id]);
  }

  Future<List<SendOrderToApiJob>> getPendingJobs() async {
    final db = await database;

    var maps = await db.query(
        EnumHelper.enumToStr(Tables.tbSendOrdersToApiJobs),
        where:
            'status = ? AND datetime(timestamp) >= datetime(\'now\', \'-12 Hour\')',
        whereArgs: ["Pending"]);

    var jobs = maps
        .map((e) => SendOrderToApiJob(
              e["id"] as int,
              DateTime.fromMicrosecondsSinceEpoch(e["time"] as int),
              e["request"] as String,
              e["response"] as String,
              e["exception"] as String,
              e["retryCount"] as int,
              e["status"] as String,
              e["authorization"] as String,
              e["idempotencyKey"] as String,
            ))
        .toList();

    return jobs;
  }

  Future<List<SendOrderToApiJob>> getFailedJobs() async {
    final db = await database;

    var maps = await db.query(
        EnumHelper.enumToStr(Tables.tbSendOrdersToApiJobs),
        where:
            'status = ? AND datetime(timestamp) < datetime(\'now\', \'-12 Hour\')',
        whereArgs: ["Pending"]);

    var jobs = maps
        .map((e) => SendOrderToApiJob(
              e["id"] as int,
              DateTime.fromMicrosecondsSinceEpoch(e["time"] as int),
              e["request"] as String,
              e["response"] as String,
              e["exception"] as String,
              e["retryCount"] as int,
              e["status"] as String,
              e["authorization"] as String,
              e["idempotencyKey"] as String,
            ))
        .toList();

    return jobs;
  }

  deleteJob(int id) async {
    final db = await database;
    db.rawDelete(
        'DELETE FROM ${EnumHelper.enumToStr(Tables.tbSendOrdersToApiJobs)} WHERE id = ?',
        [id]);
  }
}

enum Tables {
  tbSendOrdersToApiJobs,
}

class SendOrderToApiJob {
  int id;
  DateTime time;
  String request;
  String response;
  String exception;
  int retryCount;
  String status;
  String authorization;
  String idempotencyKey;

  SendOrderToApiJob(
      this.id,
      this.time,
      this.request,
      this.response,
      this.exception,
      this.retryCount,
      this.status,
      this.authorization,
      this.idempotencyKey);
}
