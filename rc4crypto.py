import math
import xlrd, xlwt
from xlutils.copy import copy

key = "_thd_salt__key_"
kb = bytes(key, encoding = "utf8")



def rc4Decode(text, key):
    res = bytearray()
    key_len = len(key)
    #1. init S-box
    box = list(range(256))                          #put 0-255 into S-box

    for i in range(256):
        box[i] = 9*i+7

    j = 0
    for i in range(256):                            #shuffle elements in S-box according to key
        # print(j, box[i])
        j = (j + box[i] + key[i % key_len]) % 256
        # print(i, j, key[i%key_len], box[i])
        box[i], box[j] = box[j], box[i]             #swap elements

    text = bytearray.fromhex(text)

    #2. make sure all elements in S-box swapped at least once
    i = j = 0
    for element in text:
        i = (i + 1) % 256
        j = (j + box[i]) % 256
        # print(element, box[(box[i] + box[j]) % 256] % 256)
        k2 = element ^ (box[(box[i] + box[j]) % 256]  % 256)
        res.append(k2)

    res[8], res[1] = res[1], res[8]
    res[9], res[3] = res[3], res[9]
    res[10], res[5] = res[5], res[10]
    res[11], res[7] = res[7], res[11]
    res = res[:8]

    res = res[::-1]                                 # little endian

    result = 0
    k = len(res) - 1
    for r in res:
        result = result + int(r) * int(math.pow(16, k*2))
        k = k - 1

    return result


def rc4decode(openid):
    openid = openid[5:]

    dataArr = openid.split("_")
    if len(dataArr) == 1:
        return openid

    return rc4Decode(dataArr[0], kb)








# styleBoldRed   = xlwt.easyxf('font: color-index red, bold on');
# headerStyle = styleBoldRed;
# wb = xlwt.Workbook()
# ws = wb.add_sheet('test')
# ws.write(0, 0, "Header",        headerStyle)
# ws.write(0, 1, "CatalogNumber", headerStyle)
# ws.write(0, 2, "PartNumber",    headerStyle)
# wb.save("D:\\test.xls")




sheetlist = ["shumei_20180723", "shumei_20180722", "shumei_20180721", "shumei_20180720", "shumei_20180719"]

excel = xlrd.open_workbook('D:\shumei_20180719-20180723.xlsx')
wb = copy(excel)

for sheet_index, sheeti in enumerate(sheetlist):
    print(sheet_index, sheeti)

    sheet = excel.sheet_by_index(sheet_index)
    w_sheet = wb.get_sheet(sheet_index)


    nrows = sheet.nrows
    ncols = sheet.ncols

    writeCol = ncols + 1

    for i in range(nrows):
        for j in range(ncols):
            if j == 1:
                openid = sheet.cell(i, j).value
                print(openid)

                ctype = 1
                xf = 0
                value = rc4decode(openid)

                # sheet.put_cell(i, writeCol, ctype, value, xf)
                w_sheet.write(i, writeCol, value)

    wb.save('D:\shumei_20180719-20180723-out.xls')