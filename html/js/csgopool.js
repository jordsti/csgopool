
function showHelp(helpId) {

    help = document.getElementById(helpId);
    
    if(help) {
      help.style.display = "inherit";
      help.style.visibility = "visible";
    }
  
}

function closeHelp(helpId) {

    help = document.getElementById(helpId);
    
    if(help) {
      help.style.display = "none";
      help.style.visibility = "hidden";
    }
}

function checkAccountForm() {
	var username = document.getElementById("username");
	var password = document.getElementById("password");
	var password2 = document.getElementById("password2");
	var email = document.getElementById("email");
	var form = document.getElementById("createaccount");
	
	if (password.value != password2.value) {
		showHelp('password_mismatch_error');
		addClass('password_group', 'has-error');
		addClass('password2_group', 'has-error');
		return;
	}
	
	if (password.value.length < passwordMinChar) {
		showHelp('password_error');
		addClass('password_group', 'has-error');
		return;
	}
	
	if (username.value.length < usernameMinChar || username.value.length > usernameMaxChar) {
		showHelp('username_error');
		addClass('username', 'has-error');
		return;
	}
	
	if (email.value.indexOf("@") != -1) {
		showHelp('email_error');
		addClass('email', 'has-error');
	}
	
	form.submit()
}

function addClass(elementId, newClass) {
	var element = document.getElementById(elementId);
	
	if(element) {
		var elementClass = element.getAttribute("class");
		element.setAttribute("class", elementClass.concat(" ", newClass));
	}
}