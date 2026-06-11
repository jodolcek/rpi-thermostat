package hr.fipu.termostat;

import android.util.Log;

import java.io.IOException;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class ApiSetPoint {
    private OkHttpClient client = new OkHttpClient();
    public void sendSetpoint(float value) {

        String json = "{ \"setpoint\": " + value + " }";

        RequestBody body = RequestBody.create(
                json,
                MediaType.get("application/json; charset=utf-8")
        );

        Request request = new Request.Builder()
                .url("http://server.apps.dj:8080/setpoint")
                .post(body)
                .build();

        client.newCall(request).enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                Log.e("API", "Failed", e);
            }

            @Override
            public void onResponse(Call call, Response response) {
                Log.d("API", "Sent setpoint: " + value);
            }
        });
    }
}
