package hr.fipu.termostat;

public class ScheduleItem {

        private String Time;
        private float SetPoint;
        ScheduleItem(String time, float setPoint) {
            this.Time = time;
            this.SetPoint = setPoint;
        }
        public String getTime() {
            return this.Time;
        }
        public float getSetPoint() {
            return this.SetPoint;
        }
    }

