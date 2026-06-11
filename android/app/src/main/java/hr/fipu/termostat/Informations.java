package hr.fipu.termostat;

public class Informations {
    private String Heating;
    private Float Temperature;
    private Float SetPoint;
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
