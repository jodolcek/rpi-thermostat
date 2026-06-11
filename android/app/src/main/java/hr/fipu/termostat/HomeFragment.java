package hr.fipu.termostat;

import android.graphics.Color;
import android.os.Bundle;

import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.TextView;

import androidx.fragment.app.Fragment;


public class HomeFragment extends Fragment {

    private TextView temp, point, heating;
    private Button plus, minus;

    private WebSocket ws;
    private ApiClient api;

    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
                             Bundle savedInstanceState) {

        View view = inflater.inflate(R.layout.fragment_home, container, false);

        temp = view.findViewById(R.id.CurrentTemp);
        point = view.findViewById(R.id.SetTempText);
        heating = view.findViewById(R.id.HeatingStatus);
        plus = view.findViewById(R.id.btnPlus);
        minus = view.findViewById(R.id.btnMinus);

        api = new ApiClient();
        ws = new WebSocket();

        ws.connect();

        ws.setCallback(info -> {
            requireActivity().runOnUiThread(() -> {
                temp.setText(info.getTemperature() + "°C");
                point.setText(info.getSetPoint() + "°C");

                if ("off".equals(info.getHeating())) {
                    heating.setText("Grijanje isključeno");
                    heating.setTextColor(Color.RED);
                } else {
                    heating.setText("Grijanje uključeno");
                    heating.setTextColor(Color.GREEN);
                }
            });
        });

        plus.setOnClickListener(v -> {
            float setPoint = Float.parseFloat(point.getText().toString().replace("°C","").trim());
            setPoint += 0.5f;
            point.setText(setPoint + "°C");
            api.sendSetpoint(setPoint);
        });

        minus.setOnClickListener(v -> {
            float setPoint = Float.parseFloat(point.getText().toString().replace("°C","").trim());
            setPoint -= 0.5f;
            point.setText(setPoint + "°C");
            api.sendSetpoint(setPoint);
        });

        return view;
    }
}