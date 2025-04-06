package io.ipsn.tordemo;
import demo.Demo;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.widget.TextView;
public class MainActivity extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        final TextView label = (TextView) findViewById(R.id.label)
        final String address = Demo.run(getFileDir().toString());
        lavel.setText(address);
    }
}
