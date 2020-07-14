/**
 *  The event handler triggered when editing the spreadsheet.
 * @param {Event} e The onEdit event.
 */
function onEdit(e) {
  // Set a comment on the edited cell to indicate when it was changed.
  //var range = e.range;
}

function onOpen(e) {
  SpreadsheetApp.getUi().createMenu("HCS Connector")
      .addItem('Open', 'showSidebar')
      .addToUi();
}

function onInstall(e) {
  onOpen(e);
}

/**
 * Opens a sidebar in the document containing the add-on's user interface.
 * This method is only used by the regular add-on, and is never called by
 * the mobile add-on version.
 */
function showSidebar() {
  var ui = HtmlService.createHtmlOutputFromFile('test-sidebar')
      .setTitle('HCS Connector');
  SpreadsheetApp.getUi().showSidebar(ui);
}

var nodeURL = "https://us-central1-baseline-spreadsheet.cloudfunctions.net"

function triggerSendProposals() {
  var response = UrlFetchApp.fetch(nodeURL + '/sheets-send-proposals');
  Logger.log(response.getContentText());
  return response.getContentText();
}