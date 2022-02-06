$(document).ready(function () {
    $("#create_subscription").click(
        function () {
            $.ajax({
                url: "/api/subscription",
                type: "POST",
                data: $('#subscription_form').serialize(),
                success: function () {
                    location.replace('/admin/');
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Incorrect Data'
                    $('#subscription_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='delete_subscription']").click(
        function () {
            var id = this.id.replace("delete_subscription_","");
            $.ajax({
                url: "/api/subscription/"+id,
                type: "DELETE",
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#subscriptions_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='enable_subscription']").click(
        function () {
            var id = this.id.replace("enable_subscription_","");
            var subscription = this.title.replace("Enable subscription ", "")
            $.ajax({
                url: "/api/subscription/"+id,
                type: "PATCH",
                data: {name: subscription, status: "on"},
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#subscriptions_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='disable_subscription']").click(
        function () {
            var id = this.id.replace("disable_subscription_","");
            var subscription = this.title.replace("Disable subscription ", "")
            $.ajax({
                url: "/api/subscription/"+id,
                type: "PATCH",
                data: {name: subscription, status: "off"},
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#subscriptions_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );
});
