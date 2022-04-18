const UPLOAD_URL = {
    "华东": "http://up.qiniu.com",
    "华南": "http://up-z2.qiniu.com",
    "华北": "http://up-z1.qiniu.com",
    "北美": "http://up-na0.qiniu.com",
    "东南亚": "http://up-as0.qiniup.com",
    "华东-浙江2": "http://up-cn-east-2.qiniup.com",
}

G.ImageBedType = "Qiniu";
G.TokenExpire = 30 * 24 * 3600 * 1000 - 10 * 1000; // Token 过期时间，单位毫秒，这表示30天
G.SortKey = "putTime"; // 用来 imageList 排序的关键字段
G.SortRule = -1; //-1表示降序排序，1表示升序排序
G.ImageKey = "key"; //七牛云返回中的 key，指定这个就不用定义 getImagesCallback

// G.getImagesCallback = function (imageList) {return imageList}; // 通过回调函数保证返回的iamgeList字段一致
// imageList = [{key: "albumName/xxxxx", putTime: 1650185929266, "fsize": , "hash": , "mimeType": }]
// key 是必须的，其他字段是可选的

G.uploadImage = function (token, file, albumName, ele, data, _success, _error) { // 必须定义上传图像函数
    let formData = new FormData();
    formData.append("token", token);
    formData.append("file", file);

    $.ajax({
        type: "POST",
        url: UPLOAD_URL[G.CloudZone],
        data: formData,
        contentType: false,
        processData: false,
        success: function (ret) {
            _success && _success(ele, data, ret, albumName);
        },
        error: function () {
            _error && _error(ele, data);
        },
    }); // end $.ajax
};

// function _uploadImage(token) {
//     // let key = G.CloudSaveKey;
//     // if (albumName.length > 0) {
//     // 	key = albumName + "/" + key;
//     // }
//     let formData = new FormData();
//     formData.append("token", token);
//     formData.append("file", file);
//     // formData.append("savekey", key);

//     $.ajax({
//         type: "POST",
//         url: G.UPLOAD_URL[G.CloudZone],
//         data: formData,
//         contentType: false,
//         processData: false,
//         success: function (ret) {
//             if (self.isAddAlbum) {
//                 self.isAddAlbum = false;
//                 $("#addOrUpdateAlbumBtn").prop('disabled', false).show();
//                 $("#albumName").val("");
//                 self.pageAddAlbum(albumName);
//                 self.toggleAddAlbum();
//                 self.images[albumName] = [];
//             }

//             data.context.remove();
//             self.addSelectedImage(G.CloudDomain + "/" + ret[IMAGEKEY]); // add image to preview
//             self.uploadRefreshImageList(albumName, ret); // reresh image list
//         },
//         error: function () {
//             data.context.empty();
//             var tpl = $('<li><div class="alert alert-danger"><a class="close" data-dismiss="alert">×</a></div></li>');
//             tpl.find('div').append('<b>Error:</b> ' + data.files[0].name + ' <small>[<i>' + formatFileSize(data.files[0].size) + '</i>]</small> ' + data.errorThrown);
//             data.context.append(tpl);
//             $("#upload-msg").scrollTop(1000);
//         }
//     });
// }