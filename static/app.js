
function setImage(loc) {
    if (loc) {
       $('#outputImage').attr('src', loc);
       $('#outputLink').attr('href', loc);
       $("#output").show();
    } else {
       $("#output").hide();
       $('#outputImage').attr('src', "");
       $('#outputLink').attr('href', "");
    }
}

function showGenerating() {
    $("#generating").show();
}

function hideGenerating() {
    $("#generating").hide();
}

function showError(err) {
    $("#error").show();
    $("#error").text(err);
}

function hideError() {
    $("#error").hide();
}

function init() {
    var frm = $('#generateForm');

    var editor = ace.edit("editor");
    editor.getSession().setMode("ace/mode/dot");



    frm.submit((function (ev) {
        $("#graphtext").val(editor.getSession().getValue());
        $.ajax({
            type: frm.attr('method'),
            url: frm.attr('action'),
            data: frm.serialize(),
            cache: false,

            beforeSend: function () {
                setImage("");
                hideError();
                showGenerating();
            },

            complete: function () {
                hideGenerating();
            },

            error: function(data) {
                hideGenerating();
                showError(data.responseText);
            },
            success: function (data) {
                hideGenerating();
                setImage(data);
            }
            });

        setImage("");
        ev.preventDefault();
    }));

    $(document).on('keydown', function(e) {
        if(e.keyCode == 13 && (e.metaKey || e.ctrlKey)) {
            frm.submit();
        }
    })
}

$(document).ready(init);

