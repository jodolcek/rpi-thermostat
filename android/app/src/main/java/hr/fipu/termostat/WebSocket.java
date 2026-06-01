
package hr.fipu.termostat;
import android.util.Log;
import org.json.JSONObject;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;
import okhttp3.WebSocketListener;

public class WebSocket {

    public interface WebSocketCallback {
        void onData(Informations informations);
    }
    private WebSocketCallback callback;

    public void setCallback(WebSocketCallback callback) {
        this.callback = callback;
    }
    public okhttp3.WebSocket webSocket;

    public void connect() {

        OkHttpClient client = new OkHttpClient();

        Request request = new Request.Builder()
                .url("ws://server.apps.dj:8080/ws")
                .build();

        webSocket = client.newWebSocket(request, new WebSocketListener() {

            @Override
            public void onOpen(okhttp3.WebSocket webSocket, Response response) {
                Log.d("WS", "Connected!");
            }
            @Override
            public void onMessage(okhttp3.WebSocket webSocket, String text) {
                try {
                    JSONObject json = new JSONObject(text);

                    float temperature = (float) json.getDouble("temperature");
                    String heating = json.getString("heating");
                    float setPoint = (float) json.getDouble("setpoint");
                   Informations informations = new Informations(heating, temperature, setPoint);

                   if (callback != null) {
                        callback.onData(informations);
                    }

                } catch (Exception e) {
                    Log.e("WS", "Parse error", e);
                }
            }

            @Override
            public void onFailure(okhttp3.WebSocket webSocket, Throwable t, Response response) {
                Log.e("WS", "Error", t);
                try {
                    Thread.sleep(5000); // čekaj 5 sekundi
                } catch (InterruptedException e) {
                    Log.e("WS", "Reconnect interrupted", e);
                }

                connect();
            }
        });
    }
}
