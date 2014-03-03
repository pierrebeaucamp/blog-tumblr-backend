import "dart:html"; 
import "dart:convert";
import 'dart:async';

List<AnchorElement> menuLinks = querySelectorAll(".link");
AnchorElement nextLink = querySelector("#next");
HtmlElement title = querySelector("title");
DivElement container = querySelector("article");
DivElement content = querySelector("main");

String last;

final NodeValidatorBuilder _htmlValidator=new NodeValidatorBuilder.common()
    ..allowElement('a', attributes: ['href']);

List<String> decodeJson(String response) {
    return JSON.decode(response);
}

void render(String id, bool push) {
    push ? container.classes.add('way-down') : null;
    if (id != "about-me" && id != "thank-you") {
        HttpRequest.getString("https://pierrebeaucamp.com/ajax/" + id).then((response) {
            var j = decodeJson(response);
            container.children.clear();
            container.nodes.add(new DocumentFragment.html('<h1 class="text-center">' + j["title"] + '</h1>'));
            container.nodes.add(new DocumentFragment.html('<p>' + j["content"] + '</p>', validator: _htmlValidator));
            container.dataset['url'] = j['url'];
            last == container.dataset['url'].split("/")[1] ? nextLink.classes.add('hidden') : nextLink.classes.remove('hidden');
            title.text = "Pierre Beaucamp | " + j["title"];
            if (push) {
                window.history.pushState(null, j["title"], j["url"] + "#nav");
                new Timer(new Duration(milliseconds: 300), () => container.classes.remove('way-down'));
            }
        });
    }
}

void menuLinkClick(Event e) {
    HtmlElement clicked = e.target;
    var id = clicked.dataset['id'];
    render(id, true);
}

void nextLinkClick(Event e) {
    String target;
    bool found = false;

    for (var m in menuLinks) {
        found ? target = m.dataset['id'] : null;
        m.dataset['id'] == container.dataset['url'].split("/")[1] ? found = true : null;
    }

    render(target, true);
}

void main() {
    for (var m in menuLinks) {
        m.href = 'javascript: void(0)';
        m.onClick.listen(menuLinkClick);
        last = m.dataset['id'];
    }

    if (container.dataset['url'].split("/")[1] != "about-me" && container.dataset['url'].split("/")[1] != "thank-you") {
        last == container.dataset['url'].split("/")[1] ? null : nextLink.classes.remove('hidden');
    }

    nextLink.href = 'javascript: void(0)';
    nextLink.onClick.listen(nextLinkClick);
    
    window.onPopState.listen((_) {
        render(window.location.pathname.split("/")[1], false);
    });
}