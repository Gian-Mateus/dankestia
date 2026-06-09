import QtQuick
import Quickshell

ShellRoot {
    id: root
    Component.onCompleted: {
        console.log("Checking Translator QML type...");
        try {
            var translator = Qt.createQmlObject('import Quickshell; Translator {}', root, "dynamic");
            console.log("Translator type exists!", translator);
        } catch (e) {
            console.log("Translator type does not exist:", e.message);
        }
        Quickshell.exit(0);
    }
}
