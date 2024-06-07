import { Toast as BootstrapToast } from 'bootstrap';
import clsx from 'clsx';

type ToastOptions = {
	title?: string;
	variant: 'info' | 'success' | 'error' | 'warning';
	body: string | Element;
	autoShow?: boolean;
	canClose?: boolean;
	container?: Element;
	additionalClasses?: string[];
} & Partial<BootstrapToast.Options>;

class Toast {
	readonly element: HTMLElement;
	private container: Element;
	private title: ToastOptions['title'];
	private body: ToastOptions['body'];
	private variant: ToastOptions['variant'];
	private canClose: ToastOptions['canClose'];
	private additionalClasses: ToastOptions['additionalClasses'];

	public instance;
	constructor(options: ToastOptions) {
		const { title, variant, autoShow = true, canClose = true, body, additionalClasses, container, ...bootstrapOptions } = options || {};
		this.container = container || document.getElementById('toastContainer')!;
		this.additionalClasses = additionalClasses;
		this.title = title || 'WowSims';
		this.variant = variant || 'info';
		this.body = body;
		this.canClose = canClose;

		this.element = this.template();
		this.container.appendChild(this.element);

		this.instance = BootstrapToast.getOrCreateInstance(this.element, {
			delay: 3000,
			...bootstrapOptions,
		});

		if (autoShow) this.instance.show();

		this.element.addEventListener('hidden.bs.toast', () => {
			this.destroy();
		});
	}

	destroy() {
		this.instance.dispose();
		this.element.remove();
	}

	show() {
		this.instance.show();
	}

	hide() {
		this.instance.hide();
	}

	getVariantIcon() {
		switch (this.variant) {
			case 'info':
				return 'fa-info-circle';
			case 'success':
				return 'fa-check-circle';
			case 'error':
				return 'fa-exclamation-circle';
			case 'warning':
				return 'fa-triangle-exclamation';
		}
	}

	template() {
		return (
			<div
				className={clsx('toast position-relative bottom-0 end-0', `toast--${this.variant}`, this.additionalClasses)}
				attributes={{
					role: 'alert',
					'aria-live': 'assertive',
					'aria-atomic': 'true',
				}}>
				<div className="toast-header">
					<i className={clsx('d-block fas fa-2xl me-2', this.getVariantIcon())}></i>
					<strong className="me-auto">{this.title}</strong>
					{this.canClose && (
						<button
							className="btn-close"
							attributes={{
								'aria-label': 'Close',
							}}
							dataset={{
								bsDismiss: 'toast',
							}}
							aria-label="Close">
							<i className={clsx('fas fa-times fa-1xl', this.getVariantIcon())}></i>
						</button>
					)}
				</div>
				<div className="toast-body">{this.body}</div>
			</div>
		) as HTMLElement;
	}
}

export default Toast;
