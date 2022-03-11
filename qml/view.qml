/****************************************************************************
**
** Copyright (C) 2015 The Qt Company Ltd.
** Contact: http://www.qt.io/licensing/
**
** This file is part of the examples of the Qt Toolkit.
**
** $QT_BEGIN_LICENSE:BSD$
** You may use this file under the terms of the BSD license as follows:
**
** "Redistribution and use in source and binary forms, with or without
** modification, are permitted provided that the following conditions are
** met:
**   * Redistributions of source code must retain the above copyright
**     notice, this list of conditions and the following disclaimer.
**   * Redistributions in binary form must reproduce the above copyright
**     notice, this list of conditions and the following disclaimer in
**     the documentation and/or other materials provided with the
**     distribution.
**   * Neither the name of The Qt Company Ltd nor the names of its
**     contributors may be used to endorse or promote products derived
**     from this software without specific prior written permission.
**
**
** THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
** "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
** LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
** A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
** OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
** SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
** LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
** DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
** THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
** (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
** OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE."
**
** $QT_END_LICENSE$
**
****************************************************************************/

//Slightly edited the original code for a scrollable TextArea and Qt Quick 2 controls

import Socks5 1.0
import QtQuick 2.2
import QtQuick.Controls 2.3
import QtQuick.Layouts 1.3
import Qt.labs.platform 1.1

ApplicationWindow {
    id:appWindow
    visible: true
    title: "代理"
    property int margin: 11
    minimumWidth: 600
    minimumHeight: 450
    maximumWidth: 600
    maximumHeight: 450
    footer:ToolBar{
        height:20
        RowLayout{
            anchors.fill:parent
            Item{
                Layout.fillWidth: true
                Text{
                    anchors.verticalCenter: parent.verticalCenter
                    text:"  qproxy "+VerGetter.VerTag+" commit: "+VerGetter.VerCommitHash+" Go "+VerGetter.GoVersion+" Qt "+VerGetter.QtVersion
                }
            }
        }
    }
    MenuBar {
        Menu {
            title: qsTr("&File")
            MenuItem {
                text: qsTr("&Quit")
                onTriggered: {
                    appTray.hide()
                    Qt.quit()
                }
            }
        }
    }
    function outputEditWithTime(msg){
        outputEdit.append(Qt.formatDateTime(new Date(),"[MM-dd hh:mm:ss.zzz]")+msg)
    }

    SystemTrayIcon {
        id:appTray
        visible:true
        icon.source:"qrc:/qml/icon.ico"
        tooltip:"qproxy:未运行"
        onActivated:function(reason){
            switch(reason){
                case SystemTrayIcon.DoubleClick:
                appWindow.show()
                appWindow.raise()
                appWindow.requestActivated()
                default:
            } 
        }
        menu:Menu {
            MenuItem {
                text: "退出"
                onTriggered: {
                    appTray.hide()
                    Qt.quit()

                }
            }
    }
    }

    Socks5server {
        id:appSocks5server
        onRunStateChange:function(isRunning){
            startButton.enabled = !isRunning
            stopButton.enabled =  isRunning
            portTextField.enabled = !isRunning
            var tip = ""
            if(isRunning){
                tip = "运行中"
            }else{
                tip = "未运行"
            }
                
            appTray.tooltip = "qproxy:"+tip
        }
        onReceiveRunningError:function(msg){
            outputEditWithTime(String(msg))
            outputEditWithTime("server stop with error")
        }
        onReceiveServingError:function(msg){
            outputEditWithTime(String(msg))
        }
    }

    ColumnLayout {
        id: mainLayout
        anchors.fill: parent
        anchors.margins: margin
        GroupBox {
            id: rowBox
            title: "socks5代理"
            Layout.fillWidth: true

            RowLayout {
                id: rowLayout
                anchors.fill: parent
                Label {
                    id: titleLabel
                    elide: Label.ElideRight
                    horizontalAlignment: Qt.AlignHCenter
                    verticalAlignment: Qt.AlignVCenter
                    Layout.fillWidth: true
                    text:"端口:"
                }
                TextField {
                    id:portTextField
                    Layout.fillWidth: true
                    validator: IntValidator {bottom: 1; top: 65535;}
                    text:"33899"
                    selectByMouse : true

                }
                Button {
                    id: startButton
                    text: "开始"
                    onClicked:function(clicked){
                        try{
                        outputEditWithTime("user start server:port="+portTextField.text)
                        appSocks5server.startServer(portTextField.text)
                        }catch(e){
                            outputEdit.append(String(e))
                        }
                    }
                }
                  Button {
                    id: stopButton
                    text: "结束"
                    enabled:false
                    onClicked:function(bool){
                        try{
                        outputEditWithTime("user stop server:port="+portTextField.text)
                        appSocks5server.stopServer()
                        }catch(e){
                        outputEdit.append(String(e))
                        }
                    }
                }
            }
        }
        GroupBox {
            id: outputBox
            title: "输出"
            implicitWidth: 200
            implicitHeight: 60
            Layout.fillWidth: true
            Layout.fillHeight: true
            Flickable {
                id: flick
                anchors.fill: parent
                contentWidth: outputEdit.paintedWidth
                contentHeight: outputEdit.paintedHeight
                clip: true

                function ensureVisible(r)
                {
                    if (contentX >= r.x)
                        contentX = r.x;
                    else if (contentX+width <= r.x+r.width)
                        contentX = r.x+r.width-width;
                    if (contentY >= r.y)
                        contentY = r.y;
                    else if (contentY+height <= r.y+r.height)
                        contentY = r.y+r.height-height;
                }

                TextEdit {
                    id: outputEdit
                    width: flick.width
                    focus: true
                    wrapMode: TextEdit.Wrap
                    onCursorRectangleChanged: flick.ensureVisible(cursorRectangle)
                    text: ""
                    selectByMouse : true
                    readOnly : true
                }
            }
            
        }
    }
}
