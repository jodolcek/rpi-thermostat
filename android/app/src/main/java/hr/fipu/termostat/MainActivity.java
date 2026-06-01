package hr.fipu.termostat;

import android.graphics.Color;
import android.os.Bundle;

import android.widget.Button;
import android.widget.TextView;

import androidx.activity.EdgeToEdge;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.graphics.Insets;
import androidx.core.view.ViewCompat;
import androidx.core.view.WindowInsetsCompat;



public class MainActivity extends AppCompatActivity {


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);


        EdgeToEdge.enable(this);
        setContentView(R.layout.activity_main);
        ViewCompat.setOnApplyWindowInsetsListener(findViewById(R.id.main), (v, insets) -> {Insets systemBars = insets.getInsets(WindowInsetsCompat.Type.systemBars());v.setPadding(systemBars.left, systemBars.top, systemBars.right, systemBars.bottom);return insets;
        });

        TextView temp, point, heating;
        Button plus, minus;
        plus = findViewById(R.id.btnPlus);
        minus = findViewById(R.id.btnMinus);

        temp = findViewById(R.id.CurrentTemp);
        point = findViewById(R.id.SetTempText);
        heating = findViewById(R.id.HeatingStatus);
        WebSocket ws = new WebSocket();
        ws.connect();
        ws.setCallback(info -> {
            runOnUiThread(() -> {
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

        plus.setOnClickListener(view -> {
            String text = point.getText().toString().replace("°C", "").trim();
            float setPoint = Float.parseFloat(text);
            setPoint += 0.5f;
            point.setText(setPoint + "°C");
        });
        minus.setOnClickListener(view -> {
            String text = point.getText().toString().replace("°C", "").trim();
            float setPoint = Float.parseFloat(text);
            setPoint -= 0.5f;
            point.setText(setPoint + "°C");
        });

    }
}