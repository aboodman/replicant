import 'dart:convert';

import 'package:flutter/services.dart';

// TODO: This string has to be statically defined on the host side, which means that
// there can only be one Replicant database per-app. Find a fix when necessary.
const channelName = 'replicant.dev';

class Replicant {
  Future<Connection> _conn;

  Replicant(String dbName) {
     _conn = Connection.create(dbName);
  }

  Future<void> putBundle(String bundle) async {
    return (await _conn)._invoke('putBundle', {'code': bundle});
  }

  // Executes the named function with provided arguments from the current
  // bundle as an atomic transaction.
  Future<dynamic> exec(String function, [List<dynamic> args = const []]) async {
    return (await _conn)._invoke('exec', {'name': function, 'args': args});
  }

  // Puts a single value into the database in its own transaction.
  Future<void> put(String id, dynamic value) async {
    return (await _conn)._invoke('put', {'id': id, 'value': value});
  }

  // Get a single value from the database.
  Future<dynamic> get(String id) async {
    return (await _conn)._invoke('get', {'id': id});
  }

  Future<void> sync(String remote) async {
    return (await _conn)._invoke("sync", {'remote': remote});
  }

  Future<void> dropDatabase() async {
    return (await _conn)._invoke('dropDatabase');
  }
}

class Connection {
  MethodChannel _platform;

  static Future<Connection> create(String dbName) async {
    var c = new Connection(new MethodChannel(channelName));
    await c._platform.invokeMethod("open", dbName);
    return c;
  }

  Connection(this._platform);

  Future<dynamic> _invoke(String name, [dynamic args = const {}]) async {
    final r = await _platform.invokeMethod(name, jsonEncode(args));
    return r == '' ? null : jsonDecode(r)['result'];
  }
}
