package hr.fipu.termostat;

import android.util.Log;

import androidx.lifecycle.LiveData;
import androidx.lifecycle.MutableLiveData;
import androidx.lifecycle.ViewModel;

import org.json.JSONObject;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import okhttp3.WebSocket;
import okhttp3.WebSocketListener;

public class WSViewModel extends ViewModel {

    private final OkHttpClient client = new OkHttpClient();

    private final MutableLiveData<Informations> information = new MutableLiveData<>();

    public LiveData<Informations> getInformations() {
        return information;
    }

    private WebSocket webSocket;

    private final ScheduledExecutorService scheduler =
            Executors.newSingleThreadScheduledExecutor();

    private boolean isReconnecting = false;
    private boolean isConnected = false;

    public void connect() {

        if (isConnected || isReconnecting) {
            Log.d("WS", "Skip connect (already active or reconnecting)");
            return;
        }

        Request request = new Request.Builder()
                .url("ws://server.apps.dj:8080/ws")
                .build();

        webSocket = client.newWebSocket(request, new WebSocketListener() {

            @Override
            public void onOpen(WebSocket webSocket, Response response) {
                Log.d("WS", "Connected");
                isConnected = true;
                isReconnecting = false;
            }

            @Override
            public void onMessage(WebSocket webSocket, String text) {
                try {
                    JSONObject json = new JSONObject(text);

                    float temperature = (float) json.getDouble("temperature");
                    String heating = json.getString("heating");
                    float setPoint = (float) json.getDouble("setpoint");
                    Informations informations = new Informations(heating, temperature, setPoint);

                    information.postValue(informations);

                } catch (Exception e) {
                    Log.e("WS", "Parse error", e);
                }
            }

            @Override
            public void onFailure(WebSocket webSocket, Throwable t, Response response) {
                Log.e("WS", "Disconnected", t);

                isConnected = false;
                scheduleReconnect();
            }
        });
    }

    private void scheduleReconnect() {

        if (isReconnecting) return;

        isReconnecting = true;

        scheduler.schedule(() -> {
            isReconnecting = false;
            connect();
        }, 5, TimeUnit.SECONDS);
    }

    @Override
    protected void onCleared() {
        super.onCleared();

        isConnected = false;
        isReconnecting = false;

        if (webSocket != null) {
            webSocket.close(1000, "ViewModel destroyed");
        }

        scheduler.shutdown();
    }
}