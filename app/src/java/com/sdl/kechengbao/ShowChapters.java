package com.sdl.kechengbao;

import android.app.Activity;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.ListView;
import android.widget.SimpleAdapter;
import org.json.JSONArray;
import org.json.JSONObject;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
*@author lincanye
*/
public class ShowChapters extends Activity {
    private String strToParse;  // 需要解析出章节信息的字符串
    private ListView lv;
    private List<Map<String, String>> mData = new ArrayList<Map<String, String>>(); // 存储的Chapter数据，用于ListView进行展示

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        Log.v("MyLog","Show Chapter OnCreate");
        /////////////////////////////////////////////
        // 从ShowCourse的bunlde获取需要解析的字符串--孙栋梁添加
        Bundle mBundle = this.getIntent().getExtras();
        Log.v("MyLog","After get bundle");
        strToParse = mBundle.getString("ChapterToParse");
        /////////////////////////////////////////////
        Log.v("MyLog","After get string:"+strToParse);
        setContentView(R.layout.activity_show_chapters);

        this.getData();   // 解析章节信息到mData中

        // 填充 ListView
        lv = (ListView)findViewById(R.id.chaptersListView);
        SimpleAdapter mSimpleAdapter = new SimpleAdapter(this,
                mData, R.layout.list_chapters, new String[]{"chapterNo", "chapterTitle", "chapterIntro", "chapterText"}, new int[]{R.id.chapterNoText, R.id.chapterTitleText, R.id.chapterIntroText, R.id.chapterText});
        lv.setAdapter(mSimpleAdapter);
    }

    /**
     * 从字符串strToParse中解析出章节信息并存放到mData中
     */
    private void getData() {
        try{
            JSONArray jsonObjs = new JSONArray(strToParse);
            Log.v("MyLog","getData: "+jsonObjs.length());
            for (int i = 0; i < jsonObjs.length(); i++) {
                JSONObject jsonObj = (JSONObject)jsonObjs.opt(i);
                String chapterNo = jsonObj.getString("No");
                String chapterTitle = jsonObj.getString("Title");
                String chapterIntro = jsonObj.getString("Intro");
                String chapterText = jsonObj.getString("Text");

                Map<String, String> mMap;
                mMap = new HashMap<String, String>();
                mMap.put("chapterNo", chapterNo);
                mMap.put("chapterTitle", chapterTitle);
                mMap.put("chapterIntro", chapterIntro);
                mMap.put("chapterText", chapterText);
                mData.add(mMap);
            }
        } catch (Exception ex) {
            ex.printStackTrace();
        }
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_faq, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }
}
