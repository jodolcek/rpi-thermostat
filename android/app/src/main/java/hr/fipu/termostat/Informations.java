package hr.fipu.termostat;

public class Informations {
    public String Heating;
    public Float Temperature;
    public Float SetPoint;
    Informations(String heating, float temperature, float setPoint) {
        this.Heating = heating;
        this.Temperature = temperature;
        this.SetPoint = setPoint;
    }
    public String getHeating() {

        return this.Heating;
    }
    public float getTemperature() {

        return this.Temperature;
    }
    public float getSetPoint() {
        return this.SetPoint;
    }

}
