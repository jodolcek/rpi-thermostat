package hr.fipu.termostat;

import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.TextView;
import androidx.fragment.app.Fragment;
import androidx.lifecycle.ViewModelProvider;
import androidx.recyclerview.widget.LinearLayoutManager;
import androidx.recyclerview.widget.RecyclerView;

public class ScheduleFragment extends Fragment {

    private ScheduleViewModel viewModel;
    private ScheduleAdapter adapter;
    private RecyclerView rv;

    @Override
    public View onCreateView(LayoutInflater inflater,
                             ViewGroup container,
                             Bundle savedInstanceState) {

        View view = inflater.inflate(R.layout.fragment_schedule, container, false);

        rv = view.findViewById(R.id.rvSchedule);

        adapter = new ScheduleAdapter();
        rv.setLayoutManager(new LinearLayoutManager(getContext()));
        rv.setAdapter(adapter);

        viewModel = new ViewModelProvider(this).get(ScheduleViewModel.class);


        viewModel.getSchedule().observe(getViewLifecycleOwner(), list -> {

            TextView empty = view.findViewById(R.id.ScEmpty);

            if (list == null || list.isEmpty()) {
                rv.setVisibility(View.GONE);
                empty.setVisibility(View.VISIBLE);
            } else {
                rv.setVisibility(View.VISIBLE);
                empty.setVisibility(View.GONE);

                adapter.setData(list);
            }
        });

        viewModel.loadSchedule();

        return view;
    }
}