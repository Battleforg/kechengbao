package com.sdl.kechengbao;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.apache.http.HttpResponse;
import org.apache.http.NameValuePair;
import org.apache.http.client.entity.UrlEncodedFormEntity;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.message.BasicNameValuePair;
import org.apache.http.params.BasicHttpParams;
import org.apache.http.params.HttpConnectionParams;
import org.apache.http.params.HttpParams;
import org.apache.http.util.EntityUtils;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;
import android.content.SharedPreferences.Editor;
import android.os.Bundle;
import android.os.Handler;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.AdapterView;
import android.widget.AdapterView.OnItemLongClickListener;
import android.widget.BaseAdapter;
import android.widget.Button;
import android.widget.ListView;
import android.widget.SimpleAdapter;

/**
*@author lincanye
*/
public class ShowNotice extends Activity {
	private List<Map<String, String>> mDataList = new ArrayList<Map<String, String>>();
	private ListView lv;
	private String courceID;
	private String userId, password,url;
	private SharedPreferences sp;
	
	public void saveSp (String key, String value) {
		Editor editor = sp.edit();
		editor.putString(key, value);
		editor.commit();
	}
	
	public JSONArray readSp (String key) {
		try {
			return new JSONArray(sp.getString(key, ""));
		} catch (JSONException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
			return null;
		}
	}
	//处理得到的Notice的JSONArray
	public void handleNotices(JSONArray notices) {
		for (int i = 0; i < notices.length(); i++) {
			JSONObject notice;
			try {
				notice = notices.getJSONObject(i);
				Map<String, String> mMap;
				mMap = new HashMap<String, String>();
				mMap.put("notice", notice.getString("Text"));
                mMap.put("timeStamp", "Posted on " + notice.getString("Date").substring(5,16));
				Log.e("mMap", mMap.toString());
				mDataList.add(mMap);
			} catch (JSONException e) {
				// TODO Auto-generated catch block
				e.printStackTrace();
			}
		}
	}
	
	public void setListView(List<Map<String, String>> nDataList) {
//	设置ListView
		lv = (ListView)findViewById(R.id.notices);
		SimpleAdapter mSimpleAdapter = new SimpleAdapter(this,
				nDataList, R.layout.list_notices, new String[]{"notice", "timeStamp"}, new int[]{R.id.listNotice, R.id.timeStampText});
		lv.setAdapter(mSimpleAdapter);
	}
	//传递消息，解析JSON并填充数据
	Handler handler = new Handler() {
		public void handleMessage(android.os.Message msg) {
			if (msg.what == 1) {
				setListView((List<Map<String, String>>)msg.obj);
			}
		}
	};
	//向服务器请求信息
	private void queryData() {
		Thread thread = new Thread(new Runnable() {
			public void run() {
				String str = "";
				JSONArray notices;
				HttpPost httppost = new HttpPost(url);
				List<NameValuePair> pair = new ArrayList<NameValuePair>();
				pair.add(new BasicNameValuePair("Action", "POLL"));
				pair.add(new BasicNameValuePair("UserID", userId));
				pair.add(new BasicNameValuePair("Password", password));
				pair.add(new BasicNameValuePair("CourseID", courceID));
                pair.add(new BasicNameValuePair("Type", "Homework"));

                Log.v("MyLog", pair.toString());

                try {
					httppost.setEntity(new UrlEncodedFormEntity(pair));
                    // 设置延时
                    HttpParams my_httpParams = new BasicHttpParams();
                    HttpConnectionParams.setConnectionTimeout(my_httpParams, 500);
                    HttpConnectionParams.setSoTimeout(my_httpParams, 500);
                    HttpResponse response = new DefaultHttpClient(my_httpParams).execute(httppost);

					str = EntityUtils.toString(response.getEntity());
//				  str = "[{\"Date\":\"2015-04-12...\",\"Text\":\"Your question got a new answer! Check it out.\"},{\"Date\":\"2015-05-14...\",\"Text\":\"TA has a new post!\"}]";
                    Log.v("MyLog", "Notice: " + str);
                    str = convertToJSONString(str);
                    Log.v("MyLog", "Notice:" + str);
                    Log.v("MyLog", "key: " + ShowNotice.this.courceID + "_Notices");
   					saveSp(ShowNotice.this.courceID + "_Notices", str);

				} catch (Exception e) {
					e.printStackTrace();
                    // TEST //
                    Log.v("MyLog","Key: "+ ShowNotice.this.courceID + "_Notices");
                    Log.v("MyLog", "Value: " + sp.getString( ShowNotice.this.courceID + "_Notices", " "));

                    if (sp.getString(ShowNotice.this.courceID + "_Notices", null) != null)
					    str = readSp(ShowNotice.this.courceID + "_Notices").toString();
                    else
                        return;
				}
				try {
					notices = new JSONArray(str);
					handleNotices(notices);
					handler.obtainMessage(1, mDataList).sendToTarget();
				} catch (JSONException e) {
					// TODO Auto-generated catch block
					e.printStackTrace();
				}
			}
		});
		thread.start();
	}

    private String convertToJSONString(String str) {
        String result = "";
        char c;
        for (int i = 0; i < str.length(); i++) {
            c = str.charAt(i);
            if (c == '\\')
                continue;

            if ((c == '\"')&&((str.charAt(i+1) == '{')||(str.charAt(i-1) == '}')))
                    continue;
            result += c;
        }

        return result;
    }

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(R.layout.activity_show_notices);
		
		sp = this.getSharedPreferences(this.userId + "_info", Context.MODE_PRIVATE);
		Bundle mBundle = this.getIntent().getExtras();
		courceID = mBundle.getString("CourseID");
        userId = mBundle.getString("UserID");
        password = mBundle.getString("Password");
        url = mBundle.getString("ServerUrl");

		queryData();
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
		if (id == R.id.action_settings) {
			return true;
		}
		return super.onOptionsItemSelected(item);
	}
}
