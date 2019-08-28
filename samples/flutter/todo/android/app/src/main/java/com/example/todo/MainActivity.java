package com.example.todo;

import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import io.flutter.app.FlutterActivity;
import io.flutter.plugin.common.MethodCall;
import io.flutter.plugin.common.MethodChannel;
import io.flutter.plugin.common.MethodChannel.MethodCallHandler;
import io.flutter.plugin.common.MethodChannel.Result;
import io.flutter.plugins.GeneratedPluginRegistrant;

import java.io.File;
import java.util.Date;

import android.util.Log;

public class MainActivity extends FlutterActivity {
  private static final String CHANNEL = "replicant.dev";

  private static repm.Connection conn;

  private static File tmpDir;

  private Handler uiThreadHandler;

  @Override
  protected void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);

    GeneratedPluginRegistrant.registerWith(this);

    uiThreadHandler = new Handler(Looper.getMainLooper());

    new MethodChannel(getFlutterView(), CHANNEL).setMethodCallHandler(
      new MethodCallHandler() {
          @Override
          public void onMethodCall(MethodCall call, Result result) {
            // TODO: Do we maybe not want to create a new thread for every call?
            // Tempting to use AsyncTask but I'm not sure how many threads the backing pool
            // has and don't want sync(), which can block for a long time, to block other
            // calls into Replicant which should be near-instant.
            new Thread(new Runnable() {
              public void run() {
                if (call.method.equals("open")) {
                  MainActivity.this.handleOpen((String)call.arguments);
                  sendResult(result, new byte[0], null);
                  return;
                }

                if (conn == null) {
                  sendResult(result, new byte[0], new Exception("Replicant database has not been opened"));
                  return;
                }

                // TODO: Avoid conversion here - can dart just send as bytes?
                byte[] argData = new byte[0];
                byte[] resultData = null;
                Exception exception = null;

                if (call.arguments != null) {
                  argData = ((String)call.arguments).getBytes();
                }

                try {
                  resultData = conn.dispatch(call.method, argData);
                } catch (Exception e) {
                  exception = e;
                }

                sendResult(result, resultData, exception);
              }
            }).start();
          }
      }
    );
  }

  private void sendResult(Result result, final byte[] data, final Exception e) {
    // TODO: Avoid conversion here - can dart accept bytes?
    final String retStr = data != null && data.length > 0 ? new String(data) : "";
    uiThreadHandler.post(new Runnable() {
      @Override
      public void run() {
        if (e != null) {
          result.error("Replicant error", e.toString(), null);
        } else {
          result.success(retStr);
        }
      }
    });
  }

  private void handleOpen(String dbName) {
    File replicantDir = new File(this.getFileStreamPath("replicant"), dbName);
    File dataDir = new File(replicantDir, "data");
    File tmpDir = new File(replicantDir, "temp");

    if (!tmpDir.exists()) {
      if (!tmpDir.mkdirs()) {
        Log.e("Replicant", "Could not create temp directory");
        return;
      }
    }
    tmpDir.deleteOnExit();

    try {
      // TODO: Properly set client ID.
      MainActivity.conn = repm.Repm.open(dataDir.getAbsolutePath(), "android/c1", tmpDir.getAbsolutePath());
    } catch (Exception e) {
      Log.e("Replicant", "Could not open Replicant database", e);
    }
  }
}
