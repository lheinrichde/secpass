// decrypt password
try {
    var crypter = sessionStorage.getItem("crypter");
    document.getElementById("passwordPre").value = sjcl.decrypt(crypter, document.getElementById("password").value);
} catch (exception) { }

// encrypt password
function manipulateForm() {
    // forward values
    document.getElementById("title").value = document.getElementById("titlePre").value;
    document.getElementById("name").value = document.getElementById("namePre").value;
    document.getElementById("mail").value = document.getElementById("mailPre").value;
    document.getElementById("url").value = document.getElementById("urlPre").value;
    document.getElementById("backupCode").value = document.getElementById("backupCodePre").value;
    document.getElementById("notes").value = document.getElementById("notesPre").value;
    document.getElementById("category").value = document.getElementById("categoryPre").value;

    // forward delete if exists
    var deletePre = document.getElementById("deletePre");
    if (deletePre != null) {
        document.getElementById("delete").checked = deletePre.checked;
    }

    // encrypt password
    var password = document.getElementById("passwordPre").value;
    var crypter = sessionStorage.getItem("crypter");
    document.getElementById("password").value = sjcl.encrypt(crypter, password, { ks: 256 });

    // submit and return
    document.getElementById("form").submit();
    return false;
}

// generate password function
function generatePassword() {
    // chars to use
    var charset = "abcdefghijklmnopqrstuvwxyz!ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_abcdefghijklmnopqrstuvwxyz#ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@";
    var password = "";

    // define password length
    var el = document.getElementById("passwordPre");
    var length = parseInt(el.value);

    // set length to 16 if not integer
    if (isNaN(length) || !isInt(el.value)) {
        length = 16;
    }

    console.log(isNaN(length) + length);

    // prevent crashing or weak passwords
    if (length > 8192) {
        length = 8192
    } else if (length < 4) {
        length = 4
    }

    // add random chars
    for (var i = 0, n = charset.length; i < length; ++i) {
        password += charset.charAt(Math.floor(Math.random() * n));
    }

    // return
    return password;
}

// check if integer
function isInt(value) {
    for (i = 0; i < value.length; i++) {
        if ((value.charAt(i) < '0') || (value.charAt(i) > '9')) return false
    }
    return true;
}

// generate random password
document.getElementById("generatePassword").addEventListener('click', function (event) {
    // get element and modify
    var el = document.getElementById("passwordPre");
    el.value = generatePassword();
});

// show password
document.getElementById("showPassword").addEventListener('click', function (event) {
    // get element
    var el = document.getElementById("passwordPre");

    // show or hide password
    if (el.type === "password") {
        el.type = "text";
    } else {
        el.type = "password";
    }
});