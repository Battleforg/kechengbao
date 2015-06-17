
jQuery(document).ready(function() {
	$("#username").focus(function(){
		$("#tip").css("display","none");
	});
	$("#password").focus(function(){
		$("#tip").css("display","none");
	});
	$("#register").click(function(){
		$("#head").text("注册");
		$("#title").text("register");
		$("#password").val("");
		$("#repassword").val("");
		$("#sure").text("注册");
		$("#register").css("display","none");
		$("#tip").css("display","none");
		$("#repassword").css("display","block");
	});
    $('.page-container form').submit(function(){
		var type = $("#sure").val();
        var username = $(this).find('.username').val();
        var password = $("#password").val();
        var repassword = $("#repassword").val();
        if(username == '') {
            $(this).find('.error').fadeOut('fast', function(){
                $(this).css('top', '27px');
            });
            $(this).find('.error').fadeIn('fast', function(){
                $(this).parent().find('.username').focus();
            });
            $("#tip").text("用户名或者密码错误");
            $("#tip").css("display","block");
            return false;
        }
        if(password == '') {
            $(this).find('.error').fadeOut('fast', function(){
                $(this).css('top', '96px');
            });
            $(this).find('.error').fadeIn('fast', function(){
                $(this).parent().find('.password').focus();
            });
            $("#tip").text("用户名或者密码错误");
            $("#tip").css("display","block");
            return false;
        }
    	if(type=="登陆"){
		    $.ajax({
				url : '172.18.35.52:8080',
				type : 'POST', 
				data : {
					'Action' : "LOGIN",
					'UserID' : username,
					'Password' : password,
				}
			 }).done(function(data) {
				alert("登陆成功");
				if(data == "200"){
					$.cookie("UserID", username, { expires: 7 });
					$.cookie("Password", password, { expires: 7 });
	//				location.href = ".html";跳转到志林界面
				}
				else{
					$("#tip").text("用户名或者密码错误");
					$("#tip").css("display","block");
					$("#password").val("");
				}
			});
    	}
    	else{
    		if(repassword==password){
    			$.ajax({
					url : '172.18.35.52:8080',
					type : 'POST',
					data : {
						'Action' : "REGISTER",
						'UserID' : username,
						'Password' : password,
					},
				}).done(function(data) {
					if(data == "200"){
						location.href="login.html";
					}
					else{
						$("#password").val("");
						$("#repassword").val("");
						$("#tip").text("注册失败：用户已存在");
						$("#tip").css("display","block");
					}
				});
		   }
    		else{
    			$("#password").val("");
				$("#repassword").val("");
    			$("#tip").text("两次密码不一致");
				$("#tip").css("display","block");
    		}
    	}
    });
    $('.page-container form .username, .page-container form .password').keyup(function(){
        $(this).parent().find('.error').fadeOut('fast');
    });
});
