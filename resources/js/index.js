// get cookie
function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');

    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];

        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }

        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }

    return "";
}

// copy password to clipboard
var copybtn = document.querySelector('.copybtn');
if (copybtn != null) {
    for (i = 0; i < copybtn.length; i++) {
        // add click listener
        copybtn[i].addEventListener('click', function (event) {
            // define textarea to copy
            var el = document.createElement("textarea");

            // get passwordid, password and cookie hash
            var passwordid = event.srcElement.getAttribute("passwordid");
            var password = document.getElementById("pw-" + passwordid).value;
            var cookie = getCookie('secpass_hash');

            // decrypt and fill element
            el.value = sjcl.decrypt(cookie, password);

            // add element to page and select
            document.body.appendChild(el);
            el.select();

            // copy and remove element
            document.execCommand('copy');
            document.body.removeChild(el);

            // hide others
            document.getElementById("passwordView").style.display = "none";
            document.getElementById("passwordEdit").style.display = "none";
            document.getElementById("passwordDelete").style.display = "none";
        });
    }
}

// show password view
var viewbtn = document.querySelectorAll('.viewbtn');
if (viewbtn != null) {
    for (i = 0; i < viewbtn.length; i++) {
        // add click listener
        viewbtn[i].addEventListener('click', function (event) {
            // get passwordid, password and cookie hash
            var passwordid = event.srcElement.getAttribute("passwordid");
            var password = document.getElementById("pw-" + passwordid).value;
            var cookie = getCookie('secpass_hash');

            // get element and decrypt
            var el = document.getElementById("passwordView");
            var decrypted = sjcl.decrypt(cookie, password);

            // check that it is not displayed
            if (el.style.display === "none" || el.value != decrypted) {
                // fill element and hide others
                el.value = decrypted;
                document.getElementById("passwordEdit").style.display = "none";
                document.getElementById("passwordDelete").style.display = "none";

                // show
                el.style.display = "block";
            } else {
                // hide
                el.style.display = "none";
            }
        });
    }
}

// show password edit form
var editbtn = document.querySelectorAll('.editbtn');
if (editbtn != null) {
    for (i = 0; i < editbtn.length; i++) {
        // add click listener
        editbtn[i].addEventListener('click', function (event) {
            // get element
            var el = document.getElementById("passwordEdit");

            // get passwordid, password and cookie hash
            var passwordid = event.srcElement.getAttribute("passwordid");
            var password = document.getElementById("pw-" + passwordid).value;
            var cookie = getCookie('secpass_hash');

            // get input elements
            var input = document.getElementById("passwordEditInput");
            var elid = document.getElementById("passwordEditID");

            // decrypt
            var decrypted = sjcl.decrypt(cookie, password);

            // check that is is not displayed
            if (el.style.display === "none" || input.value != decrypted || elid.value != passwordid) {
                // fill elements
                input.value = decrypted;
                elid.value = passwordid;

                // hide others
                document.getElementById("passwordView").style.display = "none";
                document.getElementById("passwordDelete").style.display = "none";

                // show
                el.style.display = "block";
            } else {
                // hide
                el.style.display = "none";
            }
        });
    }
}

// show password delete button
var deletebtn = document.querySelectorAll('.deletebtn');
if (deletebtn != null) {
    for (i = 0; i < deletebtn.length; i++) {
        // add click listener
        deletebtn[i].addEventListener('click', function (event) {
            // get element
            var el = document.getElementById("passwordDelete");

            // get passwordid and input element
            var passwordid = event.srcElement.getAttribute("passwordid");
            var input = document.getElementById("passwordDeleteInput");

            // check that is is not displayed
            if (el.style.display === "none" || input.value != passwordid) {
                // hide others
                document.getElementById("passwordView").style.display = "none";
                document.getElementById("passwordEdit").style.display = "none";

                // set value to passwordid and show
                input.value = passwordid;
                input.checked = false;
                el.style.display = "block";
            } else {
                //hide
                el.style.display = "none";
            }
        });
    }
}

// encrypt password
function manipulateAddPassword() {
    // forward name
    document.getElementById("name").value = document.getElementById("namePre").value;

    // encrypt password
    var password = document.getElementById("passwordPre").value;
    document.getElementById("password").value = sjcl.encrypt(getCookie('secpass_hash'), password, { ks: 256 });

    // submit and return
    document.getElementById("addPassword").submit();
    return false;
}

// encrypt password
function manipulateEditPassword() {
    // forward name
    document.getElementById("passwordEditIDAfter").value = document.getElementById("passwordEditID").value;

    // encrypt password
    var password = document.getElementById("passwordEditInput").value;
    document.getElementById("passwordEditInputAfter").value = sjcl.encrypt(getCookie('secpass_hash'), password, { ks: 256 });

    // submit and return
    document.getElementById("passwordEditAfter").submit();
    return false;
}