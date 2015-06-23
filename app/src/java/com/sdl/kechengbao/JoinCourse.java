package com.sdl.kechengbao;

import java.util.ArrayList;
import java.util.List;
import org.apache.http.NameValuePair;
import org.apache.http.client.entity.UrlEncodedFormEntity;
import org.apache.http.client.methods.HttpPost;
import org.apache.http.impl.client.DefaultHttpClient;
import org.apache.http.message.BasicNameValuePair;

import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;

/**
*@author lincanye
*/
public class JoinCourse extends Activity {
	private String courceID;
	private String userId, password;
	
	private void queryData() {
		Thread thread = new Thread(new Runnable() {
			public void run() {
				String url = "http://172.18.35.52:8080/api";
				HttpPost httppost = new HttpPost(url);
				List<NameValuePair> pair = new ArrayList<NameValuePair>();
				pair.add(new BasicNameValuePair("Action", "JOIN"));
				pair.add(new BasicNameValuePair("UserID", userId));
				pair.add(new BasicNameValuePair("Password", password));
				pair.add(new BasicNameValuePair("CourseID", courceID));
				Bundle mBundle = new Bundle();
			  mBundle.putString("CourseID", courceID);
			  mBundle.putString("UserID", userId);
			  mBundle.putString("Password", password);
			  mBundle.putString("Source", "JOIN");
			  mBundle.putString("ServerUrl", "http://172.18.33.42:8080/api");
				Intent mIntent = new Intent();
			  mIntent.setClass(JoinCourse.this, ShowCources.class);
			  mIntent.putExtras(mBundle);
				startActivity(mIntent);	
//				try {
//					httppost.setEntity(new UrlEncodedFormEntity(pair));
//					new DefaultHttpClient().execute(httppost);
//					Bundle mBundle = new Bundle();
//				  mBundle.putString("CourseID", courceID);
//				  mBundle.putString("UserID", userId);
//				  mBundle.putString("Password", password);
//				  mBundle.putString("Source", "JOIN");
//					Intent mIntent = new Intent();
//				  mIntent.setClass(JoinCourse.this, ShowCources.class);
//					mIntent.putExtras(mBundle);
//					startActivity(mIntent);	
//				} catch (Exception e) {					
//					
//					e.printStackTrace();
//				}
			}
		});
		thread.start();
	}
	
	@Override	
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		userId = "12330285";
		password = "12330";
		Bundle mBundle = this.getIntent().getExtras();
		courceID = mBundle.getString("Course");
		queryData();		
	}
}
