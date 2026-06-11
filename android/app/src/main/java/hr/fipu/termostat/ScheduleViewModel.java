package hr.fipu.termostat;

import androidx.lifecycle.LiveData;
import androidx.lifecycle.MutableLiveData;
import androidx.lifecycle.ViewModel;

import org.json.JSONArray;
import org.json.JSONObject;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;

public class ScheduleViewModel extends ViewModel {

    private MutableLiveData<List<ScheduleItem>> schedule = new MutableLiveData<>();

    public LiveData<List<ScheduleItem>> getSchedule() {
        return schedule;
    }

    public void loadSchedule() {

        OkHttpClient client = new OkHttpClient();

        Request request = new Request.Builder()
                .url("http://server.apps.dj:8080/schedule")
                .build();

        client.newCall(request).enqueue(new Callback() {

            @Override
            public void onFailure(Call call, IOException e) {
                // možeš dodati error handling
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {

                String json = response.body().string();

                List<ScheduleItem> list = new ArrayList<>();

                try {
                    JSONArray array = new JSONArray(json);

                    for (int i = 0; i < array.length(); i++) {
                        JSONObject obj = array.getJSONObject(i);

                        String time = obj.getString("time");
                        float setPoint = (float) obj.getDouble("setpoint");

                        list.add(new ScheduleItem(time, setPoint));
                    }

                } catch (Exception e) {
                    e.printStackTrace();
                }

                schedule.postValue(list);
            }
        });
    }
}
