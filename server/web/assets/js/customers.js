$(document).ready(function () {
    $("#create_customer").click(
        function () {
            $.ajax({
                url: "/api/customer",
                type: "POST",
                data: $('#customer_form').serialize(),
                success: function () {
                    location.replace('/admin/');
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Incorrect Data'
                    $('#customer_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='delete_customer']").click(
        function () {
            var id = this.id.replace("delete_customer_","");
            $.ajax({
                url: "/api/customer/"+id,
                type: "DELETE",
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#customers_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='enable_customer']").click(
        function () {
            var id = this.id.replace("enable_customer_","");
            var customer = this.title.replace("Enable subscription ", "")
            $.ajax({
                url: "/api/customer/"+id,
                type: "PATCH",
                data: {name: customer, status: "on"},
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#customers_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );

    $("a[id^='disable_customer']").click(
        function () {
            var id = this.id.replace("disable_customer_","");
            var customer = this.title.replace("Disable subscription ", "")
            $.ajax({
                url: "/api/customer/"+id,
                type: "PATCH",
                data: {name: customer, status: "off"},
                success: function () {
                    location.reload();
                },
                error: function (xhr, status, error) {
                    var err = eval("(" + xhr.responseText + ")");
                    var message = typeof err.msg !== 'undefined' && err.msg != null && err.msg != '' ? err.msg : 'Cannot delete plan'
                    $('#customers_result').html('<p>'+message+'</p>');
                }
            });
            return false;
        }
    );
});
