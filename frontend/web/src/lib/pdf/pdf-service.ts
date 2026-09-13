import { PDFDocument, rgb, StandardFonts, degrees } from 'pdf-lib';

export interface CertificateData {
	studentName: string;
	rollNumber: string;
	degree: string;
	department: string;
	issueDate: string;
	purpose: string;
	institutionName?: string;
	registrarName?: string;
}

export interface FeeReceiptData {
	receiptNumber: string;
	studentName: string;
	studentEmail: string;
	rollNumber: string;
	semester: string;
	paymentDate: string;
	paymentMethod: string;
	transactionId: string;
	items: { description: string; amount: number }[];
	notes?: string;
}

export interface HallTicketData {
	ticketNumber: string;
	studentName: string;
	rollNumber: string;
	program: string;
	semester: string;
	examCenter: string;
	exams: { courseCode: string; courseName: string; date: string; time: string; seatNo: string }[];
}

export interface TranscriptData {
	studentName: string;
	rollNumber: string;
	program: string;
	cgpa: string;
	totalCredits: number;
	semesters: {
		semesterName: string;
		gpa: string;
		courses: { code: string; title: string; credits: number; grade: string }[];
	}[];
}

export interface CustomDocumentData {
	title: string;
	subtitle?: string;
	author?: string;
	institution?: string;
	sections: { heading: string; body: string }[];
	footerText?: string;
	stamp?: 'APPROVED' | 'VERIFIED' | 'CONFIDENTIAL' | 'NONE';
}

/**
 * Creates an official University Bonafide Certificate PDF
 */
export async function generateBonafideCertificate(data: CertificateData): Promise<Uint8Array> {
	const pdfDoc = await PDFDocument.create();
	const page = pdfDoc.addPage([595.28, 841.89]); // A4 portrait
	const { width, height } = page.getSize();

	const fontBold = await pdfDoc.embedFont(StandardFonts.HelveticaBold);
	const fontRegular = await pdfDoc.embedFont(StandardFonts.Helvetica);
	const fontOblique = await pdfDoc.embedFont(StandardFonts.HelveticaOblique);

	// Border decoration
	page.drawRectangle({
		x: 24,
		y: 24,
		width: width - 48,
		height: height - 48,
		borderColor: rgb(0.12, 0.28, 0.48),
		borderWidth: 3
	});

	page.drawRectangle({
		x: 30,
		y: 30,
		width: width - 60,
		height: height - 60,
		borderColor: rgb(0.85, 0.65, 0.13),
		borderWidth: 1
	});

	// Header
	const instName = data.institutionName || 'BHAGWATI INSTITUTE OF TECHNOLOGY & SCIENCE';
	page.drawText(instName, {
		x: width / 2 - fontBold.widthOfTextAtSize(instName, 16) / 2,
		y: height - 80,
		size: 16,
		font: fontBold,
		color: rgb(0.12, 0.28, 0.48)
	});

	const subHeader = 'Affiliated to AKTU • Approved by AICTE • NAAC Accredited';
	page.drawText(subHeader, {
		x: width / 2 - fontRegular.widthOfTextAtSize(subHeader, 10) / 2,
		y: height - 100,
		size: 10,
		font: fontRegular,
		color: rgb(0.4, 0.4, 0.4)
	});

	// Certificate Title Banner
	page.drawRectangle({
		x: width / 2 - 140,
		y: height - 150,
		width: 280,
		height: 32,
		color: rgb(0.12, 0.28, 0.48)
	});

	const certTitle = 'BONAFIDE CERTIFICATE';
	page.drawText(certTitle, {
		x: width / 2 - fontBold.widthOfTextAtSize(certTitle, 14) / 2,
		y: height - 139,
		size: 14,
		font: fontBold,
		color: rgb(1, 1, 1)
	});

	// Certificate Body Text
	const dateStr = `Date of Issue: ${data.issueDate || new Date().toLocaleDateString()}`;
	page.drawText(dateStr, {
		x: 50,
		y: height - 195,
		size: 10,
		font: fontRegular,
		color: rgb(0.3, 0.3, 0.3)
	});

	const refStr = `Ref: BBDIT/ACAD/2026/CERT-${Math.floor(1000 + Math.random() * 9000)}`;
	page.drawText(refStr, {
		x: width - 50 - fontRegular.widthOfTextAtSize(refStr, 10),
		y: height - 195,
		size: 10,
		font: fontRegular,
		color: rgb(0.3, 0.3, 0.3)
	});

	const bodyP1 = `This is to certify that Mr./Ms. ${data.studentName}, Roll No: ${data.rollNumber},`;
	const bodyP2 = `is a bonafide student of this institution pursuing ${data.degree}`;
	const bodyP3 = `in the Department of ${data.department}.`;
	const bodyP4 = `This certificate is issued on the request of the student for the purpose of:`;
	const bodyP5 = `"${data.purpose || 'Official Verification & Academic Submission'}".`;

	let yPos = height - 250;
	const lines = [
		bodyP1,
		bodyP2,
		bodyP3,
		'',
		bodyP4,
		bodyP5,
		'',
		'During this tenure, their conduct and character have been found to be exemplary.'
	];

	for (const line of lines) {
		if (line) {
			page.drawText(line, {
				x: 50,
				y: yPos,
				size: 12,
				font: line.startsWith('"') ? fontOblique : fontRegular,
				color: line.startsWith('"') ? rgb(0.12, 0.28, 0.48) : rgb(0.15, 0.15, 0.15),
				lineHeight: 18
			});
		}
		yPos -= 24;
	}

	// Stamp & Signatures
	// Verified Stamp Circle
	page.drawCircle({
		x: 130,
		y: 130,
		size: 42,
		borderColor: rgb(0.12, 0.58, 0.32),
		borderWidth: 2
	});
	page.drawText('VERIFIED', {
		x: 105,
		y: 132,
		size: 10,
		font: fontBold,
		color: rgb(0.12, 0.58, 0.32)
	});
	page.drawText('OFFICIAL SEAL', {
		x: 98,
		y: 120,
		size: 8,
		font: fontBold,
		color: rgb(0.12, 0.58, 0.32)
	});

	// Signature line
	page.drawLine({
		start: { x: width - 200, y: 130 },
		end: { x: width - 50, y: 130 },
		thickness: 1,
		color: rgb(0.2, 0.2, 0.2)
	});

	page.drawText(data.registrarName || 'Dean / Registrar (Academics)', {
		x: width - 195,
		y: 112,
		size: 10,
		font: fontBold,
		color: rgb(0.15, 0.15, 0.15)
	});

	page.drawText('Authorized Institutional Signatory', {
		x: width - 195,
		y: 98,
		size: 8,
		font: fontRegular,
		color: rgb(0.4, 0.4, 0.4)
	});

	return await pdfDoc.save();
}

/**
 * Creates an official University Fee Payment Receipt PDF
 */
export async function generateFeeReceipt(data: FeeReceiptData): Promise<Uint8Array> {
	const pdfDoc = await PDFDocument.create();
	const page = pdfDoc.addPage([595.28, 841.89]);
	const { width, height } = page.getSize();

	const fontBold = await pdfDoc.embedFont(StandardFonts.HelveticaBold);
	const fontRegular = await pdfDoc.embedFont(StandardFonts.Helvetica);

	// Header background banner
	page.drawRectangle({
		x: 0,
		y: height - 100,
		width: width,
		height: 100,
		color: rgb(0.12, 0.28, 0.48)
	});

	page.drawText('CAMPUS ACADEMIC ACCOUNTS', {
		x: 40,
		y: height - 45,
		size: 16,
		font: fontBold,
		color: rgb(1, 1, 1)
	});

	page.drawText('Official University Fee Payment Receipt', {
		x: 40,
		y: height - 65,
		size: 11,
		font: fontRegular,
		color: rgb(0.85, 0.9, 0.98)
	});

	// Receipt details card
	page.drawRectangle({
		x: 40,
		y: height - 205,
		width: width - 80,
		height: 90,
		borderColor: rgb(0.85, 0.88, 0.92),
		borderWidth: 1,
		color: rgb(0.97, 0.98, 0.99)
	});

	page.drawText(`Receipt No: ${data.receiptNumber}`, { x: 55, y: height - 135, size: 10, font: fontBold });
	page.drawText(`Payment Date: ${data.paymentDate || new Date().toLocaleDateString()}`, {
		x: 55,
		y: height - 155,
		size: 10,
		font: fontRegular
	});
	page.drawText(`Transaction ID: ${data.transactionId}`, { x: 55, y: height - 175, size: 10, font: fontRegular });
	page.drawText(`Payment Mode: ${data.paymentMethod}`, { x: 55, y: height - 195, size: 10, font: fontRegular });

	page.drawText(`Student Name: ${data.studentName}`, { x: width / 2 + 20, y: height - 135, size: 10, font: fontBold });
	page.drawText(`Roll Number: ${data.rollNumber}`, { x: width / 2 + 20, y: height - 155, size: 10, font: fontRegular });
	page.drawText(`Email: ${data.studentEmail}`, { x: width / 2 + 20, y: height - 175, size: 10, font: fontRegular });
	page.drawText(`Academic Semester: ${data.semester}`, {
		x: width / 2 + 20,
		y: height - 195,
		size: 10,
		font: fontRegular
	});

	// Line items table
	const tableTop = height - 235;
	page.drawRectangle({
		x: 40,
		y: tableTop - 25,
		width: width - 80,
		height: 25,
		color: rgb(0.12, 0.28, 0.48)
	});

	page.drawText('Fee Component Description', {
		x: 55,
		y: tableTop - 18,
		size: 10,
		font: fontBold,
		color: rgb(1, 1, 1)
	});
	page.drawText('Amount (INR)', { x: width - 140, y: tableTop - 18, size: 10, font: fontBold, color: rgb(1, 1, 1) });

	let yPos = tableTop - 45;
	let totalAmount = 0;

	for (const item of data.items) {
		totalAmount += item.amount;
		page.drawText(item.description, { x: 55, y: yPos, size: 10, font: fontRegular });
		const amtStr = `INR ${item.amount.toLocaleString('en-IN')}`;
		page.drawText(amtStr, { x: width - 140, y: yPos, size: 10, font: fontBold });

		page.drawLine({
			start: { x: 40, y: yPos - 8 },
			end: { x: width - 40, y: yPos - 8 },
			thickness: 0.5,
			color: rgb(0.85, 0.85, 0.85)
		});

		yPos -= 26;
	}

	// Total Row
	page.drawRectangle({
		x: 40,
		y: yPos - 20,
		width: width - 80,
		height: 28,
		color: rgb(0.92, 0.96, 0.94)
	});

	page.drawText('TOTAL AMOUNT PAID', { x: 55, y: yPos - 12, size: 11, font: fontBold, color: rgb(0.12, 0.45, 0.25) });
	const totalStr = `INR ${totalAmount.toLocaleString('en-IN')}`;
	page.drawText(totalStr, { x: width - 150, y: yPos - 12, size: 12, font: fontBold, color: rgb(0.12, 0.45, 0.25) });

	// PAID Stamp
	page.drawRectangle({
		x: width / 2 - 60,
		y: yPos - 90,
		width: 120,
		height: 38,
		borderColor: rgb(0.12, 0.58, 0.32),
		borderWidth: 2,
		color: rgb(0.9, 0.98, 0.92)
	});
	page.drawText('PAID IN FULL', {
		x: width / 2 - 45,
		y: yPos - 75,
		size: 13,
		font: fontBold,
		color: rgb(0.12, 0.58, 0.32)
	});

	// Footer note
	page.drawText('This is a computer-generated official receipt. No physical signature is mandatory.', {
		x:
			width / 2 -
			fontRegular.widthOfTextAtSize(
				'This is a computer-generated official receipt. No physical signature is mandatory.',
				8
			) /
				2,
		y: 40,
		size: 8,
		font: fontRegular,
		color: rgb(0.5, 0.5, 0.5)
	});

	return await pdfDoc.save();
}

/**
 * Creates an Examination Hall Ticket / Admit Card PDF
 */
export async function generateHallTicket(data: HallTicketData): Promise<Uint8Array> {
	const pdfDoc = await PDFDocument.create();
	const page = pdfDoc.addPage([595.28, 841.89]);
	const { width, height } = page.getSize();

	const fontBold = await pdfDoc.embedFont(StandardFonts.HelveticaBold);
	const fontRegular = await pdfDoc.embedFont(StandardFonts.Helvetica);

	// Outer border
	page.drawRectangle({
		x: 24,
		y: 24,
		width: width - 48,
		height: height - 48,
		borderColor: rgb(0.2, 0.2, 0.2),
		borderWidth: 1.5
	});

	// Title
	const inst = 'BHAGWATI INSTITUTE OF TECHNOLOGY & SCIENCE';
	page.drawText(inst, {
		x: width / 2 - fontBold.widthOfTextAtSize(inst, 14) / 2,
		y: height - 60,
		size: 14,
		font: fontBold
	});

	const header = 'OFFICIAL EXAMINATION ADMIT CARD / HALL TICKET';
	page.drawText(header, {
		x: width / 2 - fontBold.widthOfTextAtSize(header, 11) / 2,
		y: height - 80,
		size: 11,
		font: fontBold,
		color: rgb(0.12, 0.28, 0.48)
	});

	const term = `${data.semester} • End-Semester Theory & Practical Examination`;
	page.drawText(term, {
		x: width / 2 - fontRegular.widthOfTextAtSize(term, 9) / 2,
		y: height - 96,
		size: 9,
		font: fontRegular,
		color: rgb(0.4, 0.4, 0.4)
	});

	// Candidate details box
	page.drawRectangle({
		x: 40,
		y: height - 190,
		width: width - 80,
		height: 80,
		borderColor: rgb(0.8, 0.8, 0.8),
		borderWidth: 1,
		color: rgb(0.98, 0.98, 0.98)
	});

	page.drawText(`Hall Ticket No: ${data.ticketNumber}`, { x: 55, y: height - 125, size: 10, font: fontBold });
	page.drawText(`Candidate Name: ${data.studentName}`, { x: 55, y: height - 145, size: 10, font: fontBold });
	page.drawText(`Roll Number: ${data.rollNumber}`, { x: 55, y: height - 165, size: 10, font: fontRegular });
	page.drawText(`Program: ${data.program}`, { x: 55, y: height - 180, size: 10, font: fontRegular });

	page.drawText(`Exam Center: ${data.examCenter}`, { x: width / 2 + 20, y: height - 125, size: 10, font: fontRegular });
	page.drawText('Status: ELIGIBLE', {
		x: width / 2 + 20,
		y: height - 145,
		size: 10,
		font: fontBold,
		color: rgb(0.1, 0.5, 0.2)
	});

	// Exams Table
	const tableTop = height - 220;
	page.drawRectangle({
		x: 40,
		y: tableTop - 22,
		width: width - 80,
		height: 22,
		color: rgb(0.12, 0.28, 0.48)
	});

	page.drawText('Course Code & Title', { x: 55, y: tableTop - 15, size: 9, font: fontBold, color: rgb(1, 1, 1) });
	page.drawText('Exam Date', { x: 260, y: tableTop - 15, size: 9, font: fontBold, color: rgb(1, 1, 1) });
	page.drawText('Time Slot', { x: 370, y: tableTop - 15, size: 9, font: fontBold, color: rgb(1, 1, 1) });
	page.drawText('Invigilator Sign', { x: 470, y: tableTop - 15, size: 9, font: fontBold, color: rgb(1, 1, 1) });

	let y = tableTop - 40;
	for (const ex of data.exams) {
		page.drawText(`${ex.courseCode} - ${ex.courseName}`, { x: 55, y, size: 9, font: fontRegular });
		page.drawText(ex.date, { x: 260, y, size: 9, font: fontRegular });
		page.drawText(ex.time, { x: 370, y, size: 9, font: fontRegular });

		page.drawLine({
			start: { x: 470, y: y - 2 },
			end: { x: 540, y: y - 2 },
			thickness: 0.5,
			color: rgb(0.7, 0.7, 0.7)
		});
		page.drawLine({
			start: { x: 40, y: y - 8 },
			end: { x: width - 40, y: y - 8 },
			thickness: 0.5,
			color: rgb(0.85, 0.85, 0.85)
		});

		y -= 26;
	}

	// Instructions
	page.drawText('IMPORTANT EXAMINATION INSTRUCTIONS:', {
		x: 40,
		y: y - 20,
		size: 9,
		font: fontBold,
		color: rgb(0.7, 0.1, 0.1)
	});
	const insts = [
		'1. Candidates must arrive at the examination hall 30 minutes prior to scheduled start.',
		'2. Valid institutional photo ID card must be presented along with this Hall Ticket.',
		'3. Mobile phones, smartwatches, and programmable calculators are strictly prohibited.',
		'4. Candidates will not be permitted to leave the hall before 60 minutes of commencement.'
	];

	let instY = y - 36;
	for (const i of insts) {
		page.drawText(i, { x: 40, y: instY, size: 8, font: fontRegular, color: rgb(0.3, 0.3, 0.3) });
		instY -= 14;
	}

	// Signature boxes
	page.drawLine({ start: { x: 50, y: 70 }, end: { x: 180, y: 70 }, thickness: 1, color: rgb(0.3, 0.3, 0.3) });
	page.drawText('Candidate Signature', { x: 65, y: 55, size: 9, font: fontRegular });

	page.drawLine({
		start: { x: width - 200, y: 70 },
		end: { x: width - 50, y: 70 },
		thickness: 1,
		color: rgb(0.3, 0.3, 0.3)
	});
	page.drawText('Controller of Examinations', { x: width - 185, y: 55, size: 9, font: fontBold });

	return await pdfDoc.save();
}

/**
 * Modifies an existing PDF by adding an official stamp, watermark, or signature annotation
 */
export async function annotatePdf(
	pdfBytes: ArrayBuffer | Uint8Array,
	options: {
		stampText?: string;
		watermarkText?: string;
		signatureImageDataUrl?: string;
		signaturePosition?: { x: number; y: number; pageIndex?: number };
		dateStamp?: boolean;
	}
): Promise<Uint8Array> {
	const pdfDoc = await PDFDocument.load(pdfBytes);
	const fontBold = await pdfDoc.embedFont(StandardFonts.HelveticaBold);
	const fontRegular = await pdfDoc.embedFont(StandardFonts.Helvetica);
	const pages = pdfDoc.getPages();

	if (pages.length === 0) return new Uint8Array();

	const targetPage = pages[options.signaturePosition?.pageIndex ?? 0];
	const { width, height } = targetPage.getSize();

	// Watermark across all pages if specified
	if (options.watermarkText) {
		for (const page of pages) {
			const { width: pWidth, height: pHeight } = page.getSize();
			page.drawText(options.watermarkText, {
				x: pWidth / 4,
				y: pHeight / 2 - 50,
				size: 48,
				font: fontBold,
				color: rgb(0.85, 0.85, 0.85),
				opacity: 0.35,
				rotate: degrees(45)
			});
		}
	}

	// Official Stamp Badge
	if (options.stampText) {
		targetPage.drawRectangle({
			x: width - 180,
			y: height - 80,
			width: 150,
			height: 38,
			borderColor: rgb(0.12, 0.58, 0.32),
			borderWidth: 2,
			color: rgb(0.94, 0.99, 0.95),
			opacity: 0.95
		});
		targetPage.drawText(options.stampText, {
			x: width - 170,
			y: height - 65,
			size: 11,
			font: fontBold,
			color: rgb(0.12, 0.58, 0.32)
		});
		if (options.dateStamp) {
			const dateStr = new Date().toLocaleDateString();
			targetPage.drawText(`Date: ${dateStr}`, {
				x: width - 170,
				y: height - 76,
				size: 8,
				font: fontRegular,
				color: rgb(0.12, 0.58, 0.32)
			});
		}
	}

	// Embed Signature Image
	if (options.signatureImageDataUrl) {
		const base64Data = options.signatureImageDataUrl.split(',')[1];
		if (base64Data) {
			const imageBytes = Uint8Array.from(atob(base64Data), (c) => c.charCodeAt(0));
			let embeddedImage;
			if (options.signatureImageDataUrl.includes('image/png')) {
				embeddedImage = await pdfDoc.embedPng(imageBytes);
			} else {
				embeddedImage = await pdfDoc.embedJpg(imageBytes);
			}
			const posX = options.signaturePosition?.x ?? width - 180;
			const posY = options.signaturePosition?.y ?? 80;
			targetPage.drawImage(embeddedImage, {
				x: posX,
				y: posY,
				width: 120,
				height: 50
			});
		}
	}

	return await pdfDoc.save();
}

/**
 * Creates custom Blank Document with title, author, and multi-paragraph sections
 */
export async function generateCustomDocument(data: CustomDocumentData): Promise<Uint8Array> {
	const pdfDoc = await PDFDocument.create();
	let page = pdfDoc.addPage([595.28, 841.89]);
	let { width, height } = page.getSize();

	const fontBold = await pdfDoc.embedFont(StandardFonts.HelveticaBold);
	const fontRegular = await pdfDoc.embedFont(StandardFonts.Helvetica);

	// Header
	page.drawText(data.institution || 'CAMPUS MANAGEMENT ACADEMIC NETWORK', {
		x: 50,
		y: height - 50,
		size: 10,
		font: fontBold,
		color: rgb(0.12, 0.28, 0.48)
	});

	page.drawLine({
		start: { x: 50, y: height - 58 },
		end: { x: width - 50, y: height - 58 },
		thickness: 1,
		color: rgb(0.8, 0.8, 0.8)
	});

	// Document Title
	page.drawText(data.title, {
		x: 50,
		y: height - 90,
		size: 18,
		font: fontBold,
		color: rgb(0.1, 0.1, 0.1)
	});

	if (data.subtitle) {
		page.drawText(data.subtitle, {
			x: 50,
			y: height - 110,
			size: 11,
			font: fontRegular,
			color: rgb(0.4, 0.4, 0.4)
		});
	}

	let yPos = height - 140;

	for (const sec of data.sections) {
		if (yPos < 120) {
			page = pdfDoc.addPage([595.28, 841.89]);
			yPos = height - 60;
		}

		page.drawText(sec.heading, {
			x: 50,
			y: yPos,
			size: 13,
			font: fontBold,
			color: rgb(0.12, 0.28, 0.48)
		});
		yPos -= 18;

		// Wrap body text
		const words = sec.body.split(' ');
		let currentLine = '';
		for (const word of words) {
			const testLine = currentLine ? `${currentLine} ${word}` : word;
			const lineWidth = fontRegular.widthOfTextAtSize(testLine, 10);
			if (lineWidth > width - 100) {
				page.drawText(currentLine, { x: 50, y: yPos, size: 10, font: fontRegular, color: rgb(0.2, 0.2, 0.2) });
				yPos -= 14;
				currentLine = word;
				if (yPos < 60) {
					page = pdfDoc.addPage([595.28, 841.89]);
					yPos = height - 60;
				}
			} else {
				currentLine = testLine;
			}
		}
		if (currentLine) {
			page.drawText(currentLine, { x: 50, y: yPos, size: 10, font: fontRegular, color: rgb(0.2, 0.2, 0.2) });
			yPos -= 24;
		}
	}

	return await pdfDoc.save();
}

export interface OfficialApplicationPdfData {
	templateCode: string;
	templateName: string;
	referenceNumber?: string;
	dateStr?: string;
	recipient: string;
	subject: string;
	bodyParagraphs: string[];
	closing: string;
	studentName: string;
	rollNumber: string;
	className: string;
	sessionYear?: string;
	email?: string;
	signatureDataUrl?: string;
	institutionName?: string;
}

/**
 * Creates an authentic official student application letter PDF (A4 format)
 */
export async function generateOfficialApplicationPdf(data: OfficialApplicationPdfData): Promise<Uint8Array> {
	const pdfDoc = await PDFDocument.create();
	const page = pdfDoc.addPage([595.28, 841.89]); // A4 portrait
	const { width, height } = page.getSize();

	const fontBold = await pdfDoc.embedFont(StandardFonts.HelveticaBold);
	const fontRegular = await pdfDoc.embedFont(StandardFonts.Helvetica);
	const fontOblique = await pdfDoc.embedFont(StandardFonts.HelveticaOblique);

	// Elegant Double Border
	page.drawRectangle({
		x: 30,
		y: 30,
		width: width - 60,
		height: height - 60,
		borderColor: rgb(0.12, 0.28, 0.48),
		borderWidth: 2
	});
	page.drawRectangle({
		x: 35,
		y: 35,
		width: width - 70,
		height: height - 70,
		borderColor: rgb(0.85, 0.65, 0.13),
		borderWidth: 0.8
	});

	// University Header
	const instName = data.institutionName || 'CAMPUS MANAGEMENT INSTITUTE OF HIGHER EDUCATION';
	page.drawText(instName, {
		x: width / 2 - fontBold.widthOfTextAtSize(instName, 14) / 2,
		y: height - 70,
		size: 14,
		font: fontBold,
		color: rgb(0.12, 0.28, 0.48)
	});

	const subHeader = "Accredited Grade 'A' University • Office of Academic & Student Affairs";
	page.drawText(subHeader, {
		x: width / 2 - fontRegular.widthOfTextAtSize(subHeader, 9) / 2,
		y: height - 85,
		size: 9,
		font: fontRegular,
		color: rgb(0.4, 0.4, 0.4)
	});

	// Divider line
	page.drawLine({
		start: { x: 50, y: height - 95 },
		end: { x: width - 50, y: height - 95 },
		thickness: 1,
		color: rgb(0.8, 0.8, 0.8)
	});

	// Ref and Date
	const refText = data.referenceNumber || `Ref: APP-2026-${data.templateCode.slice(0, 4)}`;
	const dateText = data.dateStr || `Date: ${new Date().toLocaleDateString('en-GB')}`;
	page.drawText(refText, {
		x: 55,
		y: height - 112,
		size: 9,
		font: fontRegular,
		color: rgb(0.3, 0.3, 0.3)
	});
	page.drawText(dateText, {
		x: width - 55 - fontRegular.widthOfTextAtSize(dateText, 9),
		y: height - 112,
		size: 9,
		font: fontRegular,
		color: rgb(0.3, 0.3, 0.3)
	});

	let currentY = height - 140;

	// Recipient Address Block
	const recipientLines = data.recipient.split('\n');
	for (const line of recipientLines) {
		page.drawText(line, {
			x: 55,
			y: currentY,
			size: 10,
			font: fontBold,
			color: rgb(0.15, 0.15, 0.15)
		});
		currentY -= 14;
	}

	currentY -= 8;

	// Subject Box
	const subjectText = `Subject: ${data.subject.replace(/^Subject:\s*/i, '')}`;
	page.drawRectangle({
		x: 50,
		y: currentY - 6,
		width: width - 100,
		height: 24,
		color: rgb(0.95, 0.96, 0.98),
		borderColor: rgb(0.8, 0.85, 0.92),
		borderWidth: 1
	});
	page.drawText(subjectText, {
		x: 58,
		y: currentY + 2,
		size: 10,
		font: fontBold,
		color: rgb(0.12, 0.28, 0.48)
	});

	currentY -= 36;

	// Salutation
	page.drawText('Respected Sir / Madam,', {
		x: 55,
		y: currentY,
		size: 10,
		font: fontRegular,
		color: rgb(0.1, 0.1, 0.1)
	});

	currentY -= 20;

	// Word wrapping text helper
	function drawWrappedText(text: string, x: number, y: number, maxWidth: number, fontSize: number, font: any): number {
		const words = text.split(' ');
		let line = '';
		let lineY = y;
		for (const word of words) {
			const testLine = line + (line ? ' ' : '') + word;
			const testWidth = font.widthOfTextAtSize(testLine, fontSize);
			if (testWidth > maxWidth && line) {
				page.drawText(line, { x, y: lineY, size: fontSize, font, color: rgb(0.15, 0.15, 0.15) });
				line = word;
				lineY -= fontSize + 5;
			} else {
				line = testLine;
			}
		}
		if (line) {
			page.drawText(line, { x, y: lineY, size: fontSize, font, color: rgb(0.15, 0.15, 0.15) });
			lineY -= fontSize + 5;
		}
		return lineY;
	}

	// Body Paragraphs
	const maxContentWidth = width - 110;
	for (const para of data.bodyParagraphs) {
		currentY = drawWrappedText(para, 55, currentY, maxContentWidth, 10, fontRegular);
		currentY -= 10;
	}

	currentY -= 10;

	// Closing
	const closingLines = data.closing.split('\n');
	for (const line of closingLines) {
		page.drawText(line, {
			x: 55,
			y: currentY,
			size: 10,
			font: fontRegular,
			color: rgb(0.15, 0.15, 0.15)
		});
		currentY -= 14;
	}

	currentY -= 5;

	// Embed Signature if present
	if (data.signatureDataUrl && data.signatureDataUrl.includes(',')) {
		try {
			const base64Data = data.signatureDataUrl.split(',')[1];
			const imageBytes = Uint8Array.from(atob(base64Data), (c) => c.charCodeAt(0));
			let embeddedImage;
			if (data.signatureDataUrl.includes('image/png')) {
				embeddedImage = await pdfDoc.embedPng(imageBytes);
			} else {
				embeddedImage = await pdfDoc.embedJpg(imageBytes);
			}
			page.drawImage(embeddedImage, {
				x: 55,
				y: currentY - 36,
				width: 110,
				height: 38
			});
			currentY -= 42;
		} catch {
			currentY -= 24;
		}
	} else {
		page.drawLine({
			start: { x: 55, y: currentY - 18 },
			end: { x: 170, y: currentY - 18 },
			thickness: 0.8,
			color: rgb(0.6, 0.6, 0.6)
		});
		page.drawText('[Student Signature]', {
			x: 55,
			y: currentY - 28,
			size: 8,
			font: fontOblique,
			color: rgb(0.5, 0.5, 0.5)
		});
		currentY -= 34;
	}

	// Student Name & Metadata
	page.drawText(data.studentName, {
		x: 55,
		y: currentY,
		size: 10,
		font: fontBold,
		color: rgb(0.12, 0.28, 0.48)
	});
	currentY -= 13;

	page.drawText(`Roll No: ${data.rollNumber}`, {
		x: 55,
		y: currentY,
		size: 9,
		font: fontRegular,
		color: rgb(0.3, 0.3, 0.3)
	});
	currentY -= 12;

	page.drawText(data.className, {
		x: 55,
		y: currentY,
		size: 9,
		font: fontRegular,
		color: rgb(0.3, 0.3, 0.3)
	});
	currentY -= 12;

	if (data.email) {
		page.drawText(`Email: ${data.email}`, {
			x: 55,
			y: currentY,
			size: 8,
			font: fontRegular,
			color: rgb(0.4, 0.4, 0.4)
		});
	}

	// Verification Footer
	page.drawLine({
		start: { x: 50, y: 55 },
		end: { x: width - 50, y: 55 },
		thickness: 0.5,
		color: rgb(0.8, 0.8, 0.8)
	});

	const footerText =
		'Official Computer-Generated Student Application • Digitally Processed by Campus Management Portal';
	page.drawText(footerText, {
		x: width / 2 - fontRegular.widthOfTextAtSize(footerText, 7.5) / 2,
		y: 42,
		size: 7.5,
		font: fontRegular,
		color: rgb(0.5, 0.5, 0.5)
	});

	return await pdfDoc.save();
}

/**
 * Triggers a browser/desktop file download of a PDF binary Uint8Array
 */
export function downloadPdfBlob(bytes: Uint8Array, filename: string = 'document.pdf') {
	const blob = new Blob([bytes as any], { type: 'application/pdf' });
	const url = URL.createObjectURL(blob);
	const a = document.createElement('a');
	a.href = url;
	a.download = filename.endsWith('.pdf') ? filename : `${filename}.pdf`;
	document.body.appendChild(a);
	a.click();
	document.body.removeChild(a);
	setTimeout(() => URL.revokeObjectURL(url), 3000);
}
