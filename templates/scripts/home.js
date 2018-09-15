$(document).ready(function () {

    DEFAULT_COOKIE_EXPIRE_TIME = 30;

    uname = '';
    session = '';
    uid = 0;
    currentVideo = null;
    listedVideos = null;

    session = getCookie('session');
    uname = getCookie('username');

    port_api = 9091
    port_scheduler = 9092
    port_streamserver = 9093
    port_web = 9094

    initPage(function () {
        if (listedVideos !== null) {
            currentVideo = listedVideos[0];
            if(listedVideos[0]['is_see']){
                selectVideoQinui(listedVideos[0]['id']);

            }else{
                selectVideo(listedVideos[0]['id']);
            }
        }

        $(".video-item").click(function () {
            var self = this.id
            var self_is_see = this.is_see
            listedVideos.forEach(function (item, index) {
                if (item['id'] === self) {
                    currentVideo = item;
                    self_is_see = item['is_see']
                    return
                }
            });


            if(self_is_see){
                selectVideoQinui(self);
            }else{
                selectVideo(self);
            }

        });

        $(".del-video-button").click(function () {
            var _JQthis = $(this)
            var id = _JQthis.attr("id")
            var str11 = _JQthis.siblings("a").find("div").eq(0).html();
            var str22 = $("<div><span>" + str11 + "</span></div>")
            str22.find("span").css({"color":"red","font-size":"200%"})
            $.confirm({
                title: '确认',
                content: '确认要删除 '+ str22.html() +' 视频? ' ,
                type: 'green',
                icon: 'glyphicon glyphicon-question-sign',
                buttons: {
                    ok: {
                        text: '确认',
                        btnClass: 'btn-primary',
                        action: function() {
                            deleteVideo(id, function (res, err) {
                                if (err !== null) {
                                    //window.alert("encounter an error when try to delete video: " + id);
                                    popupErrorMsg("encounter an error when try to delete video: " + id);
                                    return;
                                }

                                popupNotificationMsg("Successfully deleted video: " + id)
                                location.reload();
                            });
                        }
                    },
                    cancel: {
                        text: '取消',
                        btnClass: 'btn-primary'
                    }
                }
            });




        });

        $("#submit-comment").on('click', function () {
            var content = $("#comments-input").val();
            postComment(currentVideo['id'], content, function (res, err) {
                if (err !== null) {
                    popupErrorMsg("encounter and error when try to post a comment: " + content);
                    return;
                }

                if (res === "ok") {
                    popupNotificationMsg("New comment posted")
                    $("#comments-input").val("");

                    refreshComments(currentVideo['id']);
                }
            });
        });
    });

    // home page event registry
    $("#regbtn").on('click', function (e) {
        $("#regbtn").text('Loading...')
        e.preventDefault()
        registerUser(function (res, err) {
            if (err != null) {
                $('#regbtn').text("Register")
                popupErrorMsg('encounter an error, pls check your username or pwd');
                return;
            }

            var obj = JSON.parse(res);
            setCookie("session", obj["session_id"], DEFAULT_COOKIE_EXPIRE_TIME);
            setCookie("username", uname, DEFAULT_COOKIE_EXPIRE_TIME);
            $("#regsubmit").submit();
        });
    });

    $("#siginbtn").on('click', function (e) {

        $("#siginbtn").text('Loading...')
        e.preventDefault();
        signinUser(function (res, err) {
            if (err != null) {
                $('#siginbtn').text("Sign In");
                //window.alert('encounter an error, pls check your username or pwd')
                popupErrorMsg('encounter an error, pls check your username or pwd');
                return;
            }

            var obj = JSON.parse(res);
            setCookie("session", obj["session_id"], DEFAULT_COOKIE_EXPIRE_TIME);
            setCookie("username", uname, DEFAULT_COOKIE_EXPIRE_TIME);
            $("#siginsubmit").submit();
        });
    });

    $("#signinhref").on('click', function () {
        $("#regsubmit").hide();
        $("#siginsubmit").show();
    });

    $("#registerhref").on('click', function () {
        $("#regsubmit").show();
        $("#siginsubmit").hide();
    });

    // userhome event register
    $("#upload").on('click', function () {
        $("#uploadvideomodal").show();

    });


    $("#uploadform").on('submit', function (e) {
        e.preventDefault()
        var vname = $('#vname').val();
        var is_see = $('#is_see').val();

        $("#upload-submit").val("正在处理...");
        $("#upload-submit").attr("disabled", "disabled");


        if(!vname){
            vname = "v_" + Math.floor(Math.random()*1000000);
        }

        createVideo(vname, is_see, function (res, err) {
            if (err != null) {
                //window.alert('encounter an error when try to create video');
                popupErrorMsg('encounter an error when try to create video');
                return;
            }

            var obj = JSON.parse(res);
            var formData = new FormData();
            formData.append('file', $('#inputFile')[0].files[0]);
            var url = '';
            if(Number(is_see)){
                url = 'http://' + window.location.hostname + ':' + port_web + '/upload1/' + obj['id'];
            }else {
                url = 'http://' + window.location.hostname + ':' + port_web + '/upload/' + obj['id'];

            }

            $.ajax({
                url: url,
                //url:'http://127.0.0.1:8080/upload/dbibi',
                type: 'POST',
                data: formData,
                //headers: {'Access-Control-Allow-Origin': 'http://127.0.0.1:9000'},
                crossDomain: true,
                processData: false,  // tell jQuery not to process the data
                contentType: false,  // tell jQuery not to set contentType
                success: function (data) {
                    console.log(data);
                    $('#uploadvideomodal').hide();
                    location.reload();
                    //window.alert("hoa");
                },
                complete: function (xhr, textStatus) {
                    if (xhr.status === 204) {
                        window.alert("finish")
                        return;
                    }
                    if (xhr.status === 400) {
                        $("#uploadvideomodal").hide();
                        popupErrorMsg('file is too big');
                        return;
                    }
                }

            });
        });
    });

    $(".close").on('click', function () {
        $("#uploadvideomodal").hide();
    });

    $("#logout").on('click', function () {
        setCookie("session", "", -1)
        setCookie("username", "", -1)
    });


    $(".video-item").click(function () {
        var url = 'http://' + window.location.hostname + ':' + port_streamserver + '/videos/' + this.id
        var video = $("#curr-video");
        video[0].attr('src', url);
        video.load();
    });
});

function initPage(callback) {
    getUserId(function (res, err) {
        if (err != null) {
            //window.alert("Encountered error when loading user id");
            console.log("Encountered error when loading user id _初始化");

            return;
        }

        var obj = JSON.parse(res);
        uid = obj['id'];
        //window.alert(obj['id']);
        listAllVideos(function (res, err) {
            if (err != null) {
                //window.alert('encounter an error, pls check your username or pwd');
                popupErrorMsg('encounter an error, pls check your username or pwd');
                return;
            }
            var obj = JSON.parse(res);
            listedVideos = obj['videos'];

            obj['videos'].forEach(function (item, index) {
                var ele = htmlVideoListElement(item['id'], item['name'], item['display_ctime'], item['is_see']);
                $("#items").append(ele);
            });
            callback();
        });
    });
}

function setCookie(cname, cvalue, exmin) {
    var d = new Date();
    d.setTime(d.getTime() + (exmin * 60 * 1000));
    var expires = "expires=" + d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function getCookie(cname) {
    var name = cname + "=";
    var ca = document.cookie.split(';');
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

// DOM operations
function selectVideo(vid) {
    var url = 'http://' + window.location.hostname + ':' + port_web + '/videos/' + vid
    var video = $("#curr-video");
    $("#curr-video:first-child").attr('src', url);
    $("#curr-video-name").text(currentVideo['name']);
    $("#curr-video-ctime").text('Uploaded at: ' + currentVideo['display_ctime']);
    //currentVideoId = vid;
    refreshComments(vid);
}


function selectVideoQinui(vid) {
    var url = 'http://' + window.location.hostname + ':' + port_web + '/videos1/' + vid
    var video = $("#curr-video");
    $("#curr-video:first-child").attr('src', url);
    $("#curr-video-name").text(currentVideo['name']);
    $("#curr-video-ctime").text('Uploaded at: ' + currentVideo['display_ctime']);
    //currentVideoId = vid;
    refreshComments(vid);
}

function refreshComments(vid) {
    listAllComments(vid, function (res, err) {
        if (err !== null) {
            //window.alert("encounter an error when loading comments");
            popupErrorMsg('encounter an error when loading comments');
            return
        }

        var obj = JSON.parse(res);
        $("#comments-history").empty();
        if (obj['comments'] === null) {
            $("#comments-total").text('0 Comments');
        } else {
            $("#comments-total").text(obj['comments'].length + ' Comments');
        }
        if(obj['comments'] === null){

        }else{
            obj['comments'].forEach(function (item, index) {
                var ele = htmlCommentListElement(item['id'], item['author'], item['content']);
                $("#comments-history").append(ele);
            });
        }


    });
}

function popupNotificationMsg(msg) {
    var x = document.getElementById("snackbar");
    $("#snackbar").text(msg);
    x.className = "show";
    setTimeout(function () {
        x.className = x.className.replace("show", "");
    }, 2000);
}

function popupErrorMsg(msg) {
    var x = document.getElementById("errorbar");
    $("#errorbar").text(msg);
    x.className = "show";
    setTimeout(function () {
        x.className = x.className.replace("show", "");
    }, 2000);
}

function htmlCommentListElement(cid, author, content) {
    var ele = $('<div/>', {
        id: cid
    });

    ele.append(
        $('<div/>', {
            class: 'comment-author',
            text: author + ' says:'
        })
    );
    ele.append(
        $('<div/>', {
            class: 'comment',
            text: content
        })
    );

    ele.append('<hr style="height: 1px; border:none; color:#EDE3E1;background-color:#EDE3E1">');

    return ele;
}

function htmlVideoListElement(vid, name, ctime, is_see) {
    var ele = $('<a/>', {
        href: '#'
    });
    ele.append(
        $('<video/>', {
            width: '320',
            height: '240',
            poster: '/statics/img/preloader.jpg',
            controls: true
            //href: '#'
        })
    );
    ele.append(
        $('<div/>', {
            text: name
        })
    );
    ele.append(
        $('<div/>', {
            text: ctime
        })
    );


    var res = $('<div/>', {
        id: vid,
        is_see: is_see,
        class: 'video-item'
    }).append(ele);

    res.append(
        $('<button/>', {
            id: 'del-' + vid,
            type: 'button',
            class: 'del-video-button',
            text: 'Delete'
        })
    );

    res.append(
        $('<hr>', {
            size: '2'
        }).css('border-color', 'grey')
    );

    return res;
}

// Async ajax methods

// User operations
function registerUser(callback) {
    var username = $("#username").val();
    var pwd = $("#pwd").val();
    var apiUrl = window.location.hostname + ':' + port_web + '/api';

    if (username == '' || pwd == '') {
        callback(null, err);
    }

    var reqBody = {
        'user_name': username,
        'pwd': pwd
    }

    var dat = {
        'url': 'http://' + window.location.hostname + ':' + port_api + '/user',
        'method': 'POST',
        'req_body': JSON.stringify(reqBody)
    };


    $.ajax({
        url: 'http://' + window.location.hostname + ':' + port_web + '/api',
        type: 'post',
        data: JSON.stringify(dat),
        statusCode: {
            500: function () {
                callback(null, "internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of Signin");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            callback(null, "Error of register");
            return;
        }

        uname = username;
        callback(data, null);
    });
}

function signinUser(callback) {
    var username = $("#susername").val();
    var pwd = $("#spwd").val();
    var apiUrl = window.location.hostname + ':' + port_web + '/api';

    if (username == '' || pwd == '') {
        callback(null, err);
    }

    var reqBody = {
        'user_name': username,
        'pwd': pwd
    }

    var dat = {
        'url': 'http://' + window.location.hostname + ':' + port_api + '/user/' + username,
        'method': 'POST',
        'req_body': JSON.stringify(reqBody)
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':' + port_web + '/api',
        type: 'post',
        data: JSON.stringify(dat),
        statusCode: {
            500: function () {
                callback(null, "Internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of Signin");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            callback(null, "Error of Signin");
            return;
        }
        uname = username;

        callback(data, null);
    });
}

function getUserId(callback) {
    var dat = {
        'url': 'http://' + window.location.hostname + ':' + port_api + '/user/' + uname,
        'method': 'GET'
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':' + port_web + '/api',
        type: 'post',
        data: JSON.stringify(dat),
        headers: {'X-Session-Id': session},
        statusCode: {
            500: function () {
                callback(null, "Internal Error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of getUserId");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        callback(data, null);
    });
}

// Video operations
function createVideo(vname, is_see, callback) {
    var reqBody = {
        'author_id': uid,
        'name': vname.toString(),
        'is_see': Number(is_see)
    };

    var dat = {
        'url': 'http://' + window.location.hostname + ':' + port_api + '/user/' + uname + '/videos',
        'method': 'POST',
        'req_body': JSON.stringify(reqBody)
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':' + port_web + '/api',
        type: 'post',
        data: JSON.stringify(dat),
        headers: {'X-Session-Id': session},
        statusCode: {
            500: function () {
                callback(null, "Internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of Signin");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            callback(null, "Error of Signin");
            return;
        }
        callback(data, null);
    });
}

function listAllVideos(callback) {
    var dat = {
        'url': 'http://' + window.location.hostname + ':' + port_api + '/user/' + uname + '/videos',
        'method': 'GET',
        'req_body': ''
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':' + port_web + '/api',
        type: 'post',
        data: JSON.stringify(dat),
        headers: {'X-Session-Id': session},
        statusCode: {
            500: function () {
                callback(null, "Internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of Signin");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            callback(null, "Error of Signin");
            return;
        }
        callback(data, null);
    });
}

function deleteVideo(vid, callback) {
    var dat = {
        'url': 'http://' + window.location.hostname + ':' + port_api + '/user/' + uname + '/videos/' + vid,
        'method': 'DELETE',
        'req_body': ''
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':' + port_web + '/api',
        type: 'post',
        data: JSON.stringify(dat),
        headers: {'X-Session-Id': session},
        statusCode: {
            500: function () {
                callback(null, "Internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of Signin");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            callback(null, "Error of Signin");
            return;
        }
        callback(data, null);
    });
}

// Comments operations
function postComment(vid, content, callback) {
    var reqBody = {
        'author_id': uid,
        'content': content
    }


    var dat = {
        'url': 'http://' + window.location.hostname + ':' + port_api + '/videos/' + vid + '/comments',
        'method': 'POST',
        'req_body': JSON.stringify(reqBody)
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':' + port_web + '/api',
        type: 'post',
        data: JSON.stringify(dat),
        headers: {'X-Session-Id': session},
        statusCode: {
            500: function () {
                callback(null, "Internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of Signin");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            callback(null, "Error of Signin");
            return;
        }
        callback(data, null);
    });
}

function listAllComments(vid, callback) {
    var dat = {
        'url': 'http://' + window.location.hostname + ':' + port_api + '/videos/' + vid + '/comments',
        'method': 'GET',
        'req_body': ''
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':' + port_web + '/api',
        type: 'post',
        data: JSON.stringify(dat),
        headers: {'X-Session-Id': session},
        statusCode: {
            500: function () {
                callback(null, "Internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of Signin");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            callback(null, "Error of Signin");
            return;
        }
        callback(data, null);
    });
}






